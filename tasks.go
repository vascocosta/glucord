/*
 *  glucord, a simple general purpose bot for Discord.
 *  Copyright (C) 2022  Vasco Costa (gluon)
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
)

// The tskFeeds function runs in the background as a goroutine polling a collection of news feeds.
func tskFeeds(dg *discordgo.Session) {
	// Simple structure type used to send feed data to a go channel.
	// It stores a key that indexes each different feed and a value.
	// This allows the reading thread (this function) to access those two variables from the channel.
	// The key is required so that the reading thread can update the lastTime field of each feed.
	type FeedData struct {
		Key   int
		Value *gofeed.Feed
	}
	var timeFormat = "2006-01-02 15:04:05 +0000 UTC" // Time format string used by the time package.
	// Loop that runs every feedInterval seconds opening the feeds CSV file and fetching news.
	for {
		time.Sleep(time.Duration(feedInterval) * time.Second)
		//start := time.Now()
		feeds, err := readCSV(feedsFile)
		feedDataCh := make(chan FeedData)
		if err != nil {
			log.Println("tskFeeds:", err)
			continue
		}
		// Loop that spawns a goroutine worker thread per each feed source in the feeds CSV file.
		// The annonymous goroutine function accepts the k and v parameters, passed as arguments.
		// This is to avoid undesired indeterministic effects from using a closure as a goroutine.
		// The goroutine builds a Feed type by parsing the URL field for each feed in the CSV file.
		// A FeedData type is built and sent to the go channel to be received by the reading thread.
		for key, value := range feeds {
			go func(k int, v []string) {
				fp := gofeed.NewParser()
				feed, err := fp.ParseURL(v[1])
				if err != nil {
					log.Println("feed:", err)
					return
				}
				feedData := FeedData{k, feed}
				feedDataCh <- feedData
			}(key, value)
		}
		// Loop that runs a select on the go channel for as long as there's data to be read or until a timeout occurs.
		// In case feedData can be read from the communication channel, process all the feed items and show new ones.
		// In case this thread needs to wait more than 2 minutes to receive data from the goroutines a tiemout occurs.
		for {
			timeout := false
			select {
			case feedData := <-feedDataCh:
				for _, item := range feedData.Value.Items {
					// The lastTime variable keeps track of when the last feed item was retrieved.
					// If we cannot parse the time (first time) then we use timeFormat as lastTime.
					// We could use any time in the past here, but timeFormat is already available.
					lastTime, err := time.Parse(timeFormat, feeds[feedData.Key][3])
					if err != nil {
						lastTime, _ = time.Parse(timeFormat, timeFormat)
					}
					itemTime := item.PublishedParsed
					// We only want to show a feed item if itemTime > lastTime.
					// Additionally we also want to make sure the feed item is no older than 4 hours.
					// This assures only current news when restarting the bot or changing the feeds.
					if itemTime.After(lastTime) && time.Since((*itemTime)) < 8*time.Duration(hns) {
						if strings.Contains(item.Link, "?") && strings.Contains(item.Link, "&") {
							item.Link = strings.Split(item.Link, "?")[0]
						}
						dg.ChannelMessageSend(feeds[feedData.Key][2], item.Link)
						feeds[feedData.Key][3] = itemTime.String()
						writeCSV(feedsFile, feeds)
						time.Sleep(1 * time.Second)
					}
				}
			case <-time.After(60 * time.Second):
				timeout = true
			}
			if timeout {
				break
			}
		}
	}
}

// The tskEvents function runs in the background as a goroutine polling for new events.
func tskEvents(dg *discordgo.Session) {
	var announced [5]string                    // Small buffer to hold recently announced events.
	var index = 0                              // Index used to reference the buffer above.
	var timeFormat = "2006-01-02 15:04:05 UTC" // Time format string used by the time package.
	var mention string
	var image string
	do := NewDiscordOutput(dg, 0xb40000, ":alarm_clock: STARTING IN 5 MINUTES", "")
	do.Embeds = true
	// Loop that runs every minute opening the events CSV file and querying any event that starts within 5 minutes.
	for {
		time.Sleep(60 * time.Second)
		event, err := findNext("any", "any")
		if err != nil {
			log.Println("tskEvents:", err)
			continue
		}
		t, err := time.Parse(timeFormat, event[3])
		if err != nil {
			log.Println("tskEvents: Error parsing time.")
			continue
		}
		delta := time.Until(t)
		if delta.Minutes() > 5 {
			continue
		}
		// If the index becomes greather than what the buffer can hold, we reset it.
		// Otherwise we check if the announced buffer already contains the next event.
		// If it doesn't, the event is announced on the channel and added to the buffer.
		if index > 4 {
			index = 0
		} else {
			if !contains(announced[0:5], event[0]+" "+event[1]+" "+event[2]) {
				fields := []map[string]string{}
				category := map[string]string{
					"Name":  "Category:",
					"Value": event[0],
				}
				description := map[string]string{
					"Name":  "Event:",
					"Value": fmt.Sprintf("%s %s", event[1], event[2]),
				}
				fields = append(fields, category, description)
				if event[5] != "" {
					image = event[5]
				}
				if event[6] != "" {
					roles := map[string]string{
						"Name":  "Roles:",
						"Value": event[6],
					}
					fields = append(fields, roles)
					mention += event[6] + " "
				}
				dg.ChannelMessageSend(event[4], fmt.Sprintf("%sSTARTING IN 5 MINUTES: %s %s %s", mention, event[0], event[1], event[2]))
				do.Fields = &fields
				do.Image = &image
				do.Send(event[4])
				announced[index] = event[0] + " " + event[1] + " " + event[2]
				index++
			}
		}
	}
}

func tskStats(dg *discordgo.Session) {
	userCh := make(chan string)
	saveCh := make(chan string)
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		userCh <- m.Author.ID
	})
	stats, err := readCSV(statsFile)
	if err != nil {
		log.Println("tskStats:", err)
		return
	}
	timer := time.AfterFunc(300*time.Second, func() {
		saveCh <- "SAVE"
	})
	for {
		select {
		case user := <-userCh:
			found := false
			for i, v := range stats {
				if strings.EqualFold(user, v[0]) {
					old, err := strconv.Atoi(stats[i][1])
					if err != nil {
						log.Println("tskStats:", err)
						return
					}
					stats[i][1] = fmt.Sprintf("%d", old+1)
					found = true
				}
			}
			if !found {
				stats = append(stats, []string{user, "1"})
			}
		case <-saveCh:
			timer.Reset(300 * time.Second)
			err := writeCSV(statsFile, stats)
			if err != nil {
				log.Println("tskStats:", err)
			}
		}
	}
}

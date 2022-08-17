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
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// The findNext function receives a category and session and returns the chronologically next event matching that criteria.
func findNext(category string, session string) (event []string, err error) {
	var t time.Time
	var timeFormat = "2006-01-02 15:04:05 UTC"
	events, err := readCSV(eventsFile)
	if err != nil {
		return
	}
	// Loop through all events and get a parsed time for the event that matches the category and session criteria.
	// There are 3 special cases where the category and session can be set to the wildcard any in different ways.
	// Otherwise, use the default case to search for a specific category and session.
	for _, e := range events {
		switch {
		case strings.ToLower(category) == "any" && strings.ToLower(session) == "any":
			t, err = time.Parse(timeFormat, e[3])
			if err != nil {
				err = errors.New("Error parsing time.")
				return event, err
			}
		case strings.ToLower(category) != "any" && strings.ToLower(session) == "any":
			if strings.Contains(strings.ToLower(e[0]), strings.ToLower(category)) {
				t, err = time.Parse(timeFormat, e[3])
				if err != nil {
					err = errors.New("Error parsing time.")
					return event, err
				}
			}
		case strings.ToLower(category) == "any" && strings.ToLower(session) != "any":
			if strings.Contains(strings.ToLower(e[2]), strings.ToLower(session)) {
				t, err = time.Parse(timeFormat, e[3])
				if err != nil {
					err = errors.New("Error parsing time.")
					return event, err
				}
			}
		default:
			if strings.ToLower(e[0]) == strings.ToLower(category) && strings.ToLower(e[2]) == strings.ToLower(session) {
				t, err = time.Parse(timeFormat, e[3])
				if err != nil {
					err = errors.New("Error parsing time.")
					return event, err
				}
			}
		}
		// Get the time delta from now until the time of the event.
		// If delta is equal or greater than zero, this is the next event that will happen.
		delta := time.Until(t)
		if delta >= 0 {
			event = []string{e[0], e[1], e[2], e[3], e[4], e[5], e[6]}
			return event, nil
		}
	}
	err = errors.New("No event found.")
	return
}

// The help command receives a Discord session pointer, a channel and a search string.
// It then shows a compact help message listing all the possible commands of the bot.
func cmdHelp(dg *discordgo.Session, channel string, user string, search string) {
	do := NewDiscordOutput(dg, 0xb40000, "HELP", "")
	usage, err := readCSV(usageFile)
	if err != nil {
		do.Description = ":warning: Error getting usage messages."
		do.Send(channel)
		log.Println("cmdHelp:", err)
		return
	}
	users, err := readCSV(usersFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		do.Send(channel)
		log.Println("cmdHelp:", err)
		return
	}
	for _, u := range users {
		if strings.ToLower(u[0]) == strings.ToLower(user) {
			if strings.Contains(strings.ToLower(u[2]), "embeds") {
				do.Embeds = true
			}
		}
	}
	if search == "" {
		var commandList string
		for _, v := range usage {
			commandList += prefix + v[0] + "\n"
		}
		do.Description = commandList + "\n\nUse " + prefix + "help [command] to get help for a specific command."
		do.Send(channel)
	} else {
		for _, v := range usage {
			if strings.ToLower(v[0]) == strings.ToLower(search) {
				do.Description = prefix + v[1]
				do.Send(channel)
				return
			}
		}
		do.Description = ":warning: Command not found."
		do.Send(channel)
	}
}

// The next command receives a Discord session pointer, a channel, a user and an optional search string.
// It then queries the events CSV file and returns which event is happening next, showing it on the channel.
func cmdNext(dg *discordgo.Session, channel string, user string, search string) {
	var tz = "Europe/Berlin"
	var event []string
	var timeFormat = "2006-01-02 15:04:05 UTC"
	var image string
	do := NewDiscordOutput(dg, 0xb40000, "NEXT", "")
	users, err := readCSV(usersFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		do.Send(channel)
		log.Println("cmdNext:", err)
		return
	}
	for _, u := range users {
		if strings.ToLower(u[0]) == strings.ToLower(user) {
			tz = u[1]
			if strings.Contains(strings.ToLower(u[2]), "embeds") {
				do.Embeds = true
			}
		}
	}
	// Do some search string replacements in case there's actually a search argument.
	// Users use abreviated search terms, which are expanded for better database matching.
	// Retrieve the next event matching category or session criteria.
	// Else, simply retrieve the next event from any category or session type.
	if search != "" {
		switch strings.ToLower(search) {
		case "f1", "formula1", "formula 1":
			event, err = findNext("[Formula 1]", "any")
		case "f2", "formula2", "formula 2":
			event, err = findNext("[Formula 2]", "any")
		case "f3", "formula3", "formula 3":
			event, err = findNext("[Formula 3]", "any")
		case "q", "quali", "qualy", "qualifying",
			"f1 quali", "f1 qualy", "f1quali", "f1qualy", "f1 qualifying",
			"formula1 quali", "formula1 qualy", "formula1 qualifying":
			event, err = findNext("[Formula 1]", "Qualifying")
		case "r", "race", "f1 race", "f1race", "formula 1 race":
			event, err = findNext("[Formula 1]", "Race")
		case "s", "sprint", "sprint race",
			"f1 sprint", "f1sprint", "f1 sprint race",
			"formula1 sprint", "formula1 sprint race":
			event, err = findNext("[Formula 1]", "Sprint")
		default:
			event, err = findNext(search, "any")
		}
	} else {
		event, err = findNext("any", "any")
	}
	if err != nil {
		do.Description = ":warning: No event found."
		do.Send(channel)
		log.Println("cmdNext:", err)
		return
	}
	// Parse the time of the event, calculate time delta, do some formatting and finally show the results.
	// The times are localised as per the user's time zone before being shown.
	// The time delta between now and the next event uses modulo to perfectly round days, hour an minutes.
	t, err := time.Parse(timeFormat, event[3])
	if err != nil {
		do.Description = ":warning: Error parsing time."
		do.Send(channel)
		log.Println("cmdNext:", err)
		return
	}
	delta := time.Until(t)
	loc, err := time.LoadLocation(tz)
	if err != nil {
		do.Description = ":warning: Error converting time to user time zone. Using default one."
		do.Send(channel)
		log.Println("cmdNext:", err)
		return
	}
	t = t.In(loc)
	wday := t.Weekday().String()
	mday := t.Day()
	month := t.Month()
	hour := t.Hour()
	min := t.Minute()
	zone, offset := t.Zone()
	uoffset := offset / 3600
	delta = delta / 1000000000
	days := int((delta % (86400 * 30)) / 86400)
	hours := int((delta % 86400) / 3600)
	minutes := int((delta % 3600) / 60)
	fields := []map[string]string{}
	date := map[string]string{
		"Name":  "Date:",
		"Value": fmt.Sprintf("%s, %d %s", wday, mday, month),
	}
	schedule := map[string]string{
		"Name":  "Time:",
		"Value": fmt.Sprintf("%02d:%02d %s (UTC+%d)", hour, min, zone, uoffset),
	}
	category := map[string]string{
		"Name":  "Category:",
		"Value": fmt.Sprintf("%s", event[0]),
	}
	description := map[string]string{
		"Name":  "Event:",
		"Value": fmt.Sprintf("%s %s", event[1], event[2]),
	}
	countdown := map[string]string{
		"Name":  "Countdown:",
		"Value": fmt.Sprintf("%d day(s), %d hour(s), %d minute(s)", days, hours, minutes),
	}
	fields = append(fields, date, schedule, category, description, countdown)
	if event[5] != "" {
		image = event[5]
	}
	do.Color = 0x3f82ef
	do.Fields = &fields
	do.Image = &image
	do.Send(channel)

}

// The ask command receives a Discord session pointer, a channel and an arguments slice of strings.
// It then checks if the user has asked a question and displays a random answer on the channel.
func cmdAsk(dg *discordgo.Session, channel string, user string, args []string) {
	// Get a collection of answers stored as a CSV file.
	do := NewDiscordOutput(dg, 0xb40000, "ASK", "")
	users, err := readCSV(usersFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		do.Send(channel)
		log.Println("cmdAsk:", err)
		return
	}
	for _, u := range users {
		if strings.ToLower(u[0]) == strings.ToLower(user) {
			if strings.Contains(strings.ToLower(u[2]), "embeds") {
				do.Embeds = true
			}
		}
	}
	answers, err := readCSV(answersFile)
	if err != nil {
		do.Description = ":warning: Error getting answer."
		do.Send(channel)
		log.Println("cmdAsk:", err)
		return
	}
	// If the number of arguments is greater than 0, a question was asked, we show a random answer.
	// We seed the randomizer with some variable number, the current time in nano seconds.
	// Then we set the index to the answers to a random number between 0 and the length of answers.
	// Finally we show a random answer on the channel.
	if len(args) > 0 {
		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(len(answers))
		do.Color = 0x3f82ef
		do.Description = answers[index][0]
		do.Send(channel)
		// Otherwise, if we get here, it means the user didn't use the command correctly.
		// Ttherefore we show a usage message on the channel.
	} else {
		do.Description = ":warning: Usage: !ask <question>"
		do.Send(channel)
	}
}

// The plugin command receives a name, a Discord session pointer, a channel, a user and an arguments slice of strings.
// It then tries to execute the given plugin name if a file with that name is found on the plugins folder.
func cmdPlugin(name string, dg *discordgo.Session, channel string, user string, args []string) {
	var cmd *exec.Cmd
	do := NewDiscordOutput(dg, 0xb40000, strings.ToUpper(name), "")
	users, err := readCSV(usersFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		do.Send(channel)
		log.Println("cmdPlugin:", err)
		return
	}
	for _, u := range users {
		if strings.ToLower(u[0]) == strings.ToLower(user) {
			if strings.Contains(strings.ToLower(u[2]), "embeds") {
				do.Embeds = true
			}
		}
	}
	if !fileExists(pluginsFolder + name) {
		do.Description = ":warning: Unkown command or plugin."
		do.Send(channel)
		return
	}
	if len(args) == 0 {
		cmd = exec.Command(pluginsFolder+name, user)
	} else {
		var fullArgs []string
		fullArgs = append(fullArgs, user)
		fullArgs = append(fullArgs, args...)
		cmd = exec.Command(pluginsFolder+name, fullArgs...)
	}
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		do.Description = ":warning: Error executing plugin."
		do.Send(channel)
		log.Println("cmdPlugin:", err)
		return
	}
	do.Color = 0x3f82ef
	do.Description = string(cmdOutput)
	do.Send(channel)
}

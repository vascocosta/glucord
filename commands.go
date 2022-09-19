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
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	owm "github.com/briandowns/openweathermap"
	"github.com/bwmarrin/discordgo"
	"github.com/wcharczuk/go-chart"
)

// Slash commands are defined within a single var block instead of top level functions like regular commands.
// The actual execution of slash commmands is done by the equivalent regular command functions defined below.
// Slash commands are simple declarative boilerplate code to allow any regular commands to become slash ones.
var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "ask",
			Description: "Get a random answer for a question.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "question",
					Description: "The question you want to ask the bot.",
					Required:    true,
				},
			},
		},
		{
			Name:        "help",
			Description: "Show help messages for each command.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "command",
					Description: "Show the help message of a specific command.",
					Required:    false,
				},
			},
		},
		{
			Name:        "next",
			Description: "Show the next upcoming event.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "category",
					Description: "Search a specific category of events.",
					Required:    false,
				},
			},
		},
		{
			Name:        "ping",
			Description: "Send a pong in reply to a ping.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "target",
					Description: "Who to ping.",
					Required:    false,
				},
			},
		},
		{
			Name:        "poll",
			Description: "Make a channel poll.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "question",
					Description: "What to ask on this poll.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_1",
					Description: "First option.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_2",
					Description: "Second option.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_3",
					Description: "Third option.",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_4",
					Description: "Fourth option.",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_5",
					Description: "Fifth option.",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option_6",
					Description: "Sixth option.",
					Required:    false,
				},
			},
		},
		{
			Name:        "register",
			Description: "Register your user on the bot.",
		},
		{
			Name:        "roles",
			Description: "Add/remove user to/from server roles.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "role",
					Description: "The name of the role you and to add/remove.",
					Required:    false,
				},
			},
		},
		{
			Name:        "weather",
			Description: "Show the current weather for a locattion.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "location",
					Description: "Location for which you want to fetch the weather.",
					Required:    false,
				},
			},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ask": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string
			var embed *discordgo.MessageEmbed
			var embeds []*discordgo.MessageEmbed
			var args []string
			options := i.ApplicationCommandData().Options
			for _, v := range options {
				args = append(args, v.Value.(string))
			}
			do := cmdAsk(s, "", i.Member.User.ID, args)
			if do.Embeds {
				embed = do.Embed()
				embeds = append(embeds, embed)
			} else {
				content = do.Text()
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Embeds:  embeds,
				},
			})
		},
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string
			var embed *discordgo.MessageEmbed
			var embeds []*discordgo.MessageEmbed
			var args []string
			options := i.ApplicationCommandData().Options
			for _, v := range options {
				args = append(args, v.Value.(string))
			}
			do := cmdHelp(s, "", i.Member.User.ID, strings.Join(args, " "))
			if do.Embeds {
				embed = do.Embed()
				embeds = append(embeds, embed)
			} else {
				content = do.Text()
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Embeds:  embeds,
					Flags:   1 << 6,
				},
			})
		},
		"next": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string
			var embed *discordgo.MessageEmbed
			var embeds []*discordgo.MessageEmbed
			var args []string
			options := i.ApplicationCommandData().Options
			for _, v := range options {
				args = append(args, v.Value.(string))
			}
			do := cmdNext(s, "", i.Member.User.ID, strings.Join(args, " "))
			if do.Embeds {
				embed = do.Embed()
				embeds = append(embeds, embed)
			} else {
				content = do.Text()
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Embeds:  embeds,
					Flags:   1 << 6,
				},
			})
		},
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string
			var embed *discordgo.MessageEmbed
			var embeds []*discordgo.MessageEmbed
			var args []string
			options := i.ApplicationCommandData().Options
			for _, v := range options {
				args = append(args, v.Value.(string))
			}
			do := cmdPing(s, "", i.Member.User.ID, args)
			if do.Embeds {
				embed = do.Embed()
				embeds = append(embeds, embed)
			} else {
				content = do.Text()
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Embeds:  embeds,
				},
			})
		},
		"poll": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string
			var embed *discordgo.MessageEmbed
			var embeds []*discordgo.MessageEmbed
			var args []string
			options := i.ApplicationCommandData().Options
			for _, v := range options {
				args = append(args, v.Value.(string))
			}
			do := cmdPoll(s, "", i.Member.User.ID, args)
			if do.Embeds {
				embed = do.Embed()
				embeds = append(embeds, embed)
			} else {
				content = do.Text()
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Embeds:  embeds,
				},
			})
			channel, _ := s.Channel(i.ChannelID)
			pollMessageID := channel.LastMessageID
			optionsUnicode := []string{"🇦", "🇧", "🇨", "🇩", "🇪", "🇫"}
			for i := 0; i != len(options)-1; i++ {
				s.MessageReactionAdd(channel.ID, pollMessageID, optionsUnicode[i])
			}
			go func() {
				time.Sleep(5 * time.Minute)
				results := make(map[string]int)
				for i, v := range optionsUnicode {
					users, err := s.MessageReactions(channel.ID, pollMessageID, v, 0, "", "")
					if err != nil {
						log.Println(err)
					}
					if len(users) > 0 {
						results[fmt.Sprintf("%s - %s", v, args[i+1])] = len(users) - 1
					}
				}
				scoreList := make(ScoreList, len(results))
				i := 0
				for k, v := range results {
					scoreList[i] = Score{k, v}
					i++
				}
				sort.Sort(sort.Reverse(scoreList))
				s.ChannelMessageSend(channel.ID, "The poll has ended, here are the results:\n")
				for _, v := range scoreList {
					s.ChannelMessageSend(channel.ID, fmt.Sprintf("%s: %d votes\n", v.Key, v.Points))
				}
				err := s.MessageReactionsRemoveAll(channel.ID, pollMessageID)
				if err != nil {
					log.Println(err)
				}
			}()
		},
		"register": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string
			var embed *discordgo.MessageEmbed
			var embeds []*discordgo.MessageEmbed
			do := cmdRegister(s, "", i.Member.User.ID)
			if do.Embeds {
				embed = do.Embed()
				embeds = append(embeds, embed)
			} else {
				content = do.Text()
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Embeds:  embeds,
					Flags:   1 << 6,
				},
			})
		},
		"roles": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string
			var embed *discordgo.MessageEmbed
			var embeds []*discordgo.MessageEmbed
			var args []string
			options := i.ApplicationCommandData().Options
			for _, v := range options {
				args = append(args, v.Value.(string))
			}
			do := cmdRoles(s, "", i.Member.User.ID, args)
			if do.Embeds {
				embed = do.Embed()
				embeds = append(embeds, embed)
			} else {
				content = do.Text()
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Embeds:  embeds,
				},
			})
		},
		"weather": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string
			var embed *discordgo.MessageEmbed
			var embeds []*discordgo.MessageEmbed
			var args []string
			options := i.ApplicationCommandData().Options
			for _, v := range options {
				args = append(args, v.Value.(string))
			}
			do := cmdWeather(s, "", i.Member.User.ID, args)
			if do.Embeds {
				embed = do.Embed()
				embeds = append(embeds, embed)
			} else {
				content = do.Text()
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
					Embeds:  embeds,
				},
			})
		},
	}
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
				err = errors.New("error parsing time")
				return event, err
			}
		case strings.ToLower(category) != "any" && strings.ToLower(session) == "any":
			if strings.Contains(strings.ToLower(e[0]), strings.ToLower(category)) {
				t, err = time.Parse(timeFormat, e[3])
				if err != nil {
					err = errors.New("error parsing time")
					return event, err
				}
			}
		case strings.ToLower(category) == "any" && strings.ToLower(session) != "any":
			if strings.Contains(strings.ToLower(e[2]), strings.ToLower(session)) {
				t, err = time.Parse(timeFormat, e[3])
				if err != nil {
					err = errors.New("error parsing time")
					return event, err
				}
			}
		default:
			if strings.EqualFold(e[0], category) && strings.EqualFold(e[2], session) {
				t, err = time.Parse(timeFormat, e[3])
				if err != nil {
					err = errors.New("error parsing time")
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
	err = errors.New("no event found")
	return
}

// The ask command receives a Discord session pointer, a channel and an arguments slice of strings.
// It then checks if the user has asked a question and displays a random answer on the channel.
func cmdAsk(dg *discordgo.Session, channel string, user string, args []string) (do *DiscordOutput) {
	do = NewDiscordOutput(dg, 0xb40000, "ASK", "")
	users, err := readCSV(usersFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		log.Println("cmdAsk:", err)
		return
	}
	for _, u := range users {
		if strings.EqualFold(u[0], user) {
			if strings.Contains(strings.ToLower(u[2]), "embeds") {
				do.Embeds = true
			}
		}
	}
	// Get a collection of answers stored as a CSV file.
	answers, err := readCSV(answersFile)
	if err != nil {
		do.Description = ":warning: Error getting answer."
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
		do.Description = fmt.Sprintf("**Question:** %s\n\n**Answer:** %s", strings.Join(args, " "), answers[index][0])
		// Otherwise, if we get here, it means the user didn't use the command correctly.
		// Ttherefore we show a usage message on the channel.
	} else {
		do.Description = ":warning: Usage: !ask <question>"
	}
	return
}

// The bet command receives a Discord session pointer, a channel, a user and a bet containing 3 drivers.
// It then stores the bet provided by the user, or lets the user know his current bet for the next race.
func cmdBet(dg *discordgo.Session, channel string, user string, bet []string) (do *DiscordOutput) {
	var correct int
	var bets [][]string
	var update bool
	do = NewDiscordOutput(dg, 0xb40000, "BET", "")
	users, err := readCSV(usersFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		log.Println("cmdAsk:", err)
		return
	}
	for _, u := range users {
		if strings.EqualFold(u[0], user) {
			if strings.Contains(strings.ToLower(u[2]), "embeds") {
				do.Embeds = true
			}
		}
	}
	users, err = readCSV(betFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		log.Println("cmdBet:", err)
		return
	}
	if !isUser(user, users) {
		users = append(users, []string{strings.ToLower(user), "Europe/Berlin", ""})
		err = writeCSV(betFile, users)
		if err != nil {
			do.Description = ":warning: Error registering user to the bet command."
			log.Println("cmdBet:", err)
			return
		}
	}
	event, err := findNext("[formula 1]", "race")
	if err != nil {
		do.Description = ":warning: Bets are closed."
		log.Println("cmdBet:", err)
		return
	}
	bets, err = readCSV(betsFile)
	if err != nil {
		do.Description = ":warning: Error getting bets."
		log.Println("cmdBet:", err)
		return
	}
	// If no bet is provided as argument, we simply show the user's current bet, if he's placed one.
	if len(bet) == 0 {
		for i := len(bets) - 1; i >= 0; i-- {
			if strings.ToLower(bets[i][0]) == strings.ToLower(event[1]) && strings.ToLower(bets[i][1]) == strings.ToLower(user) {
				first := strings.ToUpper(bets[i][2])
				second := strings.ToUpper(bets[i][3])
				third := strings.ToUpper(bets[i][4])
				do.Description = fmt.Sprintf("Your current bet for the %s: %s %s %s", event[1], first, second, third)
				return
			}
		}
		do.Description = fmt.Sprintf("You haven't placed a bet for the %s yet.\nUse !bet log to check older bets.", event[1])
		return
	}
	drivers, err := readCSV(driversFile)
	if err != nil {
		do.Description = ":warning: Error getting drivers."
		log.Println("cmdBet:", err)
		return
	}
	// If instead of a normal bet the user provides a single word, we interpret it as an argument.
	// There are multiple arguments, for which we show the driver odds or a log of the user's bets.
	// Other possible argument is the nick of a registered user, for which we show that user's bet.
	// Alternatively, if the argument provided isn't a valid command option, we let the user know.
	if len(bet) == 1 {
		switch strings.ToLower(bet[0]) {
		case "multipliers", "odds":
			var output string
			odds, err := toStringMap(drivers, 1, 2)
			if err != nil {
				do.Description = ":warning: Error getting odds."
				log.Println("cmdBet:", err)
				return
			}
			scoreList := make(ScoreList, len(drivers))
			i := 0
			for k, v := range odds {
				integerOdds, err := strconv.Atoi(v)
				if err != nil {
					do.Description = ":warning: Error getting odds."
					log.Println("cmdBet:", err)
					return
				}
				scoreList[i] = Score{k, integerOdds}
				i++
			}
			sort.Sort(scoreList)
			for _, v := range scoreList {
				output += fmt.Sprintf("%s %d\n", strings.ToUpper(v.Key), v.Points)
			}
			if len(output) > 3 {
				do.Description = output
			}
		case "log":
			var betsFound bool
			var counter int
			for i := len(bets) - 1; i >= 0 && counter < 3; i-- {
				if strings.ToLower(bets[i][1]) == strings.ToLower(user) {
					betsFound = true
					do.Description =
						fmt.Sprintf("Your bet for the %s: %s %s %s %s points.",
							bets[i][0],
							strings.ToUpper(bets[i][2]),
							strings.ToUpper(bets[i][3]),
							strings.ToUpper(bets[i][4]),
							bets[i][5])
					counter += 1
				}
			}
			if !betsFound {
				do.Description = ":warning: No recent bets from you."
			}
		case "points":
			scoreList := make(ScoreList, len(users))
			for _, user := range users {
				points, _ := strconv.Atoi(user[2])
				if points > 0 {
					member, err := dg.GuildMember(guild, user[0])
					if err != nil {
						log.Println("cmdBet:", err)
						continue
					}
					scoreList = append(scoreList, Score{member.User.Username, points})
				}
			}
			sort.Sort(sort.Reverse(scoreList))
			output := ""
			for i, score := range scoreList {
				if score.Points > 0 {
					output += fmt.Sprintf("%d. %s %d\n", i+1, score.Key, score.Points)
				}
			}
			if len(output) > 3 {
				do.Description = output
			}
		default:
			do.Description = ":warning: Unknown command option."
		}
		return
	}
	if len(bet) != 3 {
		do.Description = ":warning: The bet must contain 3 drivers."
		return
	}
	// Finally, if we reach this point, it means the user has provided a valid bet composed of 3 drivers.
	// We verify that all 3 driver codes are valid as per the drivers CSV file before we go any further.
	// If the 3 codes are valid, we either place a new bet or update an already placed bet for the race.
	first := strings.ToLower(bet[0])
	second := strings.ToLower(bet[1])
	third := strings.ToLower(bet[2])
	for _, driver := range drivers {
		code := strings.ToLower(driver[1])
		if code == first || code == second || code == third {
			correct++
		}
	}
	if correct != 3 {
		do.Description = ":warning: Invalid drivers."
		return
	}
	for i := 0; i < len(bets); i++ {
		if strings.ToLower(bets[i][0]) == strings.ToLower(event[1]) && strings.ToLower(bets[i][1]) == strings.ToLower(user) {
			update = true
			bets[i] = []string{event[1], strings.ToLower(user), first, second, third, "0"}
			break
		}
	}
	if !update {
		bets = append(bets, []string{event[1], strings.ToLower(user), first, second, third, "0"})
	}
	err = writeCSV(betsFile, bets)
	if err != nil {
		do.Description = ":warning: Error updating bet."
		log.Println("cmdBet:", err)
		return
	}
	do.Color = 0x3f82ef
	do.Description = "Your bet for the " + event[1] + " was successfully updated."
	return
}

// The help command receives a Discord session pointer, a channel and a search string.
// It then shows a compact help message listing all the possible commands of the bot.
func cmdHelp(dg *discordgo.Session, channel string, user string, search string) (do *DiscordOutput) {
	do = NewDiscordOutput(dg, 0xb40000, "HELP", "")
	users, err := readCSV(usersFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		log.Println("cmdHelp:", err)
		return
	}
	for _, u := range users {
		if strings.EqualFold(u[0], user) {
			if strings.Contains(strings.ToLower(u[2]), "embeds") {
				do.Embeds = true
			}
		}
	}
	// Get a collection of usage strings stored as a CSV file.
	usage, err := readCSV(usageFile)
	if err != nil {
		do.Description = ":warning: Error getting usage messages."
		log.Println("cmdHelp:", err)
		return
	}
	if search == "" {
		var commandList string
		for _, v := range usage {
			commandList += prefix + v[0] + "\n"
		}
		do.Description = commandList + "\n\nUse " + prefix + "help [command] to get help for a specific command."
	} else {
		for _, v := range usage {
			if strings.EqualFold(v[0], search) {
				do.Description = prefix + v[1]
				return
			}
		}
		do.Description = ":warning: Command not found."
	}
	return
}

// The next command receives a Discord session pointer, a channel, a user and an optional search string.
// It then queries the events CSV file and returns which event is happening next, showing it on the channel.
func cmdNext(dg *discordgo.Session, channel string, user string, search string) (do *DiscordOutput) {
	var tz = "Europe/Berlin"
	var event []string
	var timeFormat = "2006-01-02 15:04:05 UTC"
	var image string
	do = NewDiscordOutput(dg, 0xb40000, "NEXT", "")
	users, err := readCSV(usersFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		log.Println("cmdNext:", err)
		return
	}
	for _, u := range users {
		if strings.EqualFold(u[0], user) {
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
		result := ""
		result, err = lookupAlias(search)
		if err != nil {
			event, err = findNext(search, "any")
		} else {
			event, err = findNext(result, "any")
		}
	} else {
		event, err = findNext("any", "any")
	}
	if err != nil {
		do.Description = ":warning: No event found."
		return
	}
	// Parse the time of the event, calculate time delta, do some formatting and finally show the results.
	// The times are localised as per the user's time zone before being shown.
	// The time delta between now and the next event uses modulo to perfectly round days, hour an minutes.
	t, err := time.Parse(timeFormat, event[3])
	if err != nil {
		do.Description = ":warning: Error parsing time."
		log.Println("cmdNext:", err)
		return
	}
	delta := time.Until(t)
	loc, err := time.LoadLocation(tz)
	if err != nil {
		do.Description = ":warning: Error converting time to user time zone. Using default one."
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
		"Value": event[0],
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
	return
}

// The ping command receives a Discord session pointer, a channel, a user and an arguments slice of strings.
// It then answers to the user using the Pong word or the target word passed by the user as an argument.
func cmdPing(dg *discordgo.Session, channel string, user string, args []string) (do *DiscordOutput) {
	do = NewDiscordOutput(dg, 0xb40000, "PING", "")
	users, err := readCSV(usersFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		log.Println("cmdQuote:", err)
		return
	}
	for _, u := range users {
		if strings.EqualFold(u[0], user) {
			if strings.Contains(strings.ToLower(u[2]), "embeds") {
				do.Embeds = true
			}
		}
	}
	do.Color = 0x3f82ef
	// Distinguish between sending simply the word Pong or whatver word was passed as argument by the user.
	if len(args) > 0 {
		do.Description = args[0]
	} else {
		do.Description = "Pong."
	}
	return
}

// The plugin command receives a name, a Discord session pointer, a channel, a user and an arguments slice of strings.
// It then tries to execute the given plugin name if a file with that name is found on the plugins folder.
func cmdPlugin(name string, dg *discordgo.Session, channel string, user string, args []string, finishedCh chan bool) {
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
		if strings.EqualFold(u[0], user) {
			if strings.Contains(strings.ToLower(u[2]), "embeds") {
				do.Embeds = true
			}
		}
	}
	// We check if the command is a plugin or not by checking if a file with that name exists.
	// If it doesn't exist this isn't a valid plugin and therefore we must stop the execution.
	if !fileExists(pluginsFolder + name) {
		select {
		case finishedCh <- true:
			// The plugin doesn't exist, but we still send true to the finished channel.
		case <-time.After(1 * time.Second):
			// If the main thread doesn't read the channel, then timeout after 1 second.
		}
		do.Description = ":warning: Unkown command or plugin."
		do.Send(channel)
		return
	}
	// Otherwise this is a valid plugin and we execute the process with the correct arguments.
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
		select {
		case finishedCh <- true:
			// The plugin had a problem, but we still send true to the finished channel.
		case <-time.After(1 * time.Second):
			// If the main thread doesn't read the channel, then timeout after 1 second.
		}
		do.Description = ":warning: Error executing plugin."
		do.Send(channel)
		log.Println("cmdPlugin:", err)
		return
	}
	do.Color = 0x3f82ef
	do.Description = string(cmdOutput)
	split := strings.Split(string(cmdOutput), "\n")
	if len(split) > 0 {
		if strings.HasPrefix(split[0], "GLUCORD-PLUGIN-HEADER:") {
			do.Description = strings.Join(split[1:], "\n")
			if strings.Contains(split[0], "EMBEDS=OFF") {
				do.Embeds = false
			}
		}
	}
	select {
	case finishedCh <- true:
		// The plugin finished with success so we send true to the finished channel.
	case <-time.After(1 * time.Second):
		// If the main thread doesn't read the channel, then timeout after 1 second.
	}
	do.Send(channel)
}

// The poll command receives a Discord session pointer, a channel, a user and an arguments slice of strings.
// It then makes a poll on a Discord channel using the poll question and all the possible answer options.
// It then waits for votes from the users and finally displays the results of the poll after a timeout.
func cmdPoll(dg *discordgo.Session, channel string, user string, args []string) (do *DiscordOutput) {
	do = NewDiscordOutput(dg, 0xb40000, "POLL (5 min)", "")
	users, err := readCSV(usersFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		do.Send(channel)
		log.Println("cmdPoll:", err)
		return
	}
	for _, u := range users {
		if strings.EqualFold(u[0], user) {
			if strings.Contains(strings.ToLower(u[2]), "embeds") {
				do.Embeds = true
			}
		}
	}
	optionsUnicode := []string{"🇦", "🇧", "🇨", "🇩", "🇪", "🇫"}
	optionsValue := ""
	for i := 1; i != len(args); i++ {
		optionsValue += fmt.Sprintf("%s - %s\n", optionsUnicode[i-1], args[i])
	}
	fields := []map[string]string{}
	question := map[string]string{
		"Name":  "Question:",
		"Value": args[0],
	}
	options := map[string]string{
		"Name":  "Answers:",
		"Value": optionsValue,
	}
	fields = append(fields, question, options)
	do.Fields = &fields
	do.Color = 0x3f82ef
	return
}

// The processbets command receives a Discord session pointer, a channel and a nick.
// It then processes the placed bets, according to the results in the results file.
func cmdProcessBets(dg *discordgo.Session, channel string, user string) (do *DiscordOutput) {
	do = NewDiscordOutput(dg, 0xb40000, "PROCESSBETS", "")
	users, err := readCSV(usersFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		log.Println("cmdHelp:", err)
		return
	}
	for _, u := range users {
		if strings.EqualFold(u[0], user) {
			if strings.Contains(strings.ToLower(u[2]), "embeds") {
				do.Embeds = true
			}
		}
	}
	if strings.ToLower(user) != strings.ToLower("541209780929167400") {
		do.Description = ":warning: Only gluon can use this command."
		return
	}
	users, err = readCSV(betFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		log.Println("cmdProcessBets:", err)
		return
	}
	results, err := readCSV(resultsFile)
	if err != nil {
		do.Description = ":warning: Error getting results."
		log.Println("cmdProcessBets:", err)
		return
	}
	if results[0][0] == results[0][4] {
		do.Description = ":warning: " + results[0][0] + " bets have already been processed in the past."
		return
	}
	bets, err := readCSV(betsFile)
	if err != nil {
		do.Description = ":warning: Error getting bets."
		log.Println("cmdProcessBets:", err)
		return
	}
	users, err = readCSV(betFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		log.Println("cmdProcessBets:", err)
		return
	}
	drivers, err := readCSV(driversFile)
	if err != nil {
		do.Description = ":warning: Error getting drivers."
		log.Println("cmdProcessBets:", err)
		return
	}
	odds, err := toStringMap(drivers, 1, 2)
	if err != nil {
		do.Description = ":warning: Error getting odds."
		log.Println("cmdProcessBets:", err)
		return
	}
	// This is the main loop where we go through each bet placed by the user and process it.
	// If the race on the bet matches the race on the results file, we calculate its score.
	for i, bet := range bets {
		score := 0
		first := strings.ToLower(results[0][1])
		second := strings.ToLower(results[0][2])
		third := strings.ToLower(results[0][3])
		if strings.ToLower(bet[0]) == strings.ToLower(results[0][0]) {
			// If the first driver is on the podium, we have two different possibilities.
			// If the first driver is the first on the results, we score 10 * multiplier.
			// If the first driver is not the first on the results, we score 5 * multiplier.
			if contains([]string{first, second, third}, strings.ToLower(bet[2])) {
				multiplier, err := strconv.Atoi(odds[bet[2]])
				if err != nil {
					do.Description = ":warning: Error applying multiplier."
					log.Println("cmdProcessBets:", err)
					return
				}
				if strings.ToLower(bet[2]) == strings.ToLower(results[0][1]) {
					score += (10 * multiplier)
				} else {
					score += (5 * multiplier)
				}
			}
			// If the second driver is on the podium, we have two different possibilities.
			// If the second driver is the second on the results, we score 10 * multiplier.
			// If the second driver is not the second on the results, we score 5 * multiplier.
			if contains([]string{first, second, third}, strings.ToLower(bet[3])) {
				multiplier, err := strconv.Atoi(odds[bet[3]])
				if err != nil {
					do.Description = ":warning: Error applying multiplier."
					log.Println("cmdProcessBets:", err)
					return
				}
				if strings.ToLower(bet[3]) == strings.ToLower(results[0][2]) {
					score += (10 * multiplier)
				} else {
					score += (5 * multiplier)
				}
			}
			// If the third driver is on the podium, we have two different possibilities.
			// If the third driver is the third on the results, we score 10 * multiplier.
			// If the third driver is not the third on the results, we score 5 * multiplier.
			if contains([]string{first, second, third}, strings.ToLower(bet[4])) {
				multiplier, err := strconv.Atoi(odds[bet[4]])
				if err != nil {
					do.Description = ":warning: Error applying multiplier."
					log.Println("cmdProcessBets:", err)
					return
				}
				if strings.ToLower(bet[4]) == strings.ToLower(results[0][3]) {
					score += (10 * multiplier)
				} else {
					score += (5 * multiplier)
				}
			}
			bets[i][5] = strconv.Itoa(score)
			// Update the total number of points for each driver on the users file.
			// The code above only handles points for each bet, not for each user.
			for j, user := range users {
				if strings.ToLower(user[0]) == strings.ToLower(bet[1]) {
					currentScore, err := strconv.Atoi(users[j][2])
					if err != nil {
						do.Description = ":warning: Error getting current score."
						log.Println("cmdProcessBets:", err)
						return
					}
					users[j][2] = strconv.Itoa(currentScore + score)
				}
			}
		}
		err = writeCSV(betFile, users)
		if err != nil {
			do.Description = ":warning: Error storing user points."
			log.Println("cmdProcessBets:", err)
			return
		}
	}
	// Finally update the bets file with the points for each bet for the current race.
	// The results file is updated so that the last field is set to the current race.
	err = writeCSV(betsFile, bets)
	if err != nil {
		do.Description = ":warning: Error storing bet points."
		log.Println("cmdProcessBets:", err)
		return
	}
	results[0][4] = results[0][0]
	err = writeCSV(resultsFile, results)
	if err != nil {
		do.Description = ":warning: Error storing last processed bet.."
		log.Println("cmdProcessBets:", err)
		return
	}
	do.Color = 0x3f82ef
	do.Description = results[0][0] + " bets successfully processed."
	return
}

// The quote command receives a Discord session pointer, a channel and an arguments slice of strings.
// It then checks if there are arguments and displays a random quote or adds a new quote accordingly.
func cmdQuote(dg *discordgo.Session, channel string, user string, args []string) (do *DiscordOutput) {
	do = NewDiscordOutput(dg, 0xb40000, "QUOTE", "")
	users, err := readCSV(usersFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		log.Println("cmdQuote:", err)
		return
	}
	for _, u := range users {
		if strings.EqualFold(u[0], user) {
			if strings.Contains(strings.ToLower(u[2]), "embeds") {
				do.Embeds = true
			}
		}
	}
	// Get a collection of quotes stored as a CSV file.
	quotes, err := readCSV(quotesFile)
	if err != nil {
		do.Description = ":warning: Error getting quote."
		log.Println("cmdQuote:", err)
		return
	}
	// Filter only the quotes of the current channel.
	var channelQuotes [][]string
	for _, quote := range quotes {
		if strings.EqualFold(quote[2], channel) {
			channelQuotes = append(channelQuotes, quote)
		}
	}
	// If there are no arguments or if the first argument is "get", show a random quote.
	// We seed the randomizer with some variable number, the current time in nano seconds.
	// Then we set the index to the quotes to a random number between 0 and the length of quotes.
	// Finally we show a random quote on the channel.
	if len(args) == 0 || (len(args) > 0 && strings.ToLower(args[0]) == "get") {
		if len(channelQuotes) == 0 {
			do.Description = ":warning: There are no quotes for this channel."
			return
		}
		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(len(channelQuotes))
		do.Color = 0x3f82ef
		do.Description = fmt.Sprintf("%s - %s", channelQuotes[index][1], channelQuotes[index][0])
		// If there is more than one argument and the first argument is "add", add the provided quote.
		// Finally we show a confirmation message on the channel.
	} else if len(args) > 1 && strings.ToLower(args[0]) == "add" {
		quotes = append(quotes, []string{time.Now().Format("02-01-2006"), strings.Join(args[1:], " "), strings.ToLower(channel)})
		err = writeCSV(quotesFile, quotes)
		if err != nil {
			do.Description = "Error adding quote."
			log.Println("cmdQuote:", err)
			return
		}
		do.Color = 0x3f82ef
		do.Description = "Quote added."
		// Otherwise, if we get here, it means the user didn't use the command correctly.
		// Ttherefore we show a usage message on the channel.
	} else {
		do.Description = "Usage: !quote [get|add] [text]"
	}
	return
}

// The register command receives a Discord session pointer, a channel and a user.
// It then checks if the user isn't already registered and registers it with the bot.
func cmdRegister(dg *discordgo.Session, channel string, user string) (do *DiscordOutput) {
	do = NewDiscordOutput(dg, 0xb40000, "REGISTER", "")
	users, err := readCSV(usersFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		log.Println("cmdRegister:", err)
		return
	}
	for _, u := range users {
		if strings.EqualFold(u[0], user) {
			if strings.Contains(strings.ToLower(u[2]), "embeds") {
				do.Embeds = true
			}
		}
	}
	// If the user is already a known user to the bot, we don't register it.
	// Otherwise we add this new user as a registered user on the users file.
	if isUser(strings.ToLower(user), users) {
		do.Description = ":warning: You are already registered."
		return
	}
	users = append(users, []string{strings.ToLower(user), "Europe/Berlin", "embeds", ""})
	err = writeCSV(usersFile, users)
	if err != nil {
		do.Description = ":warning: Error registering user."
		log.Println("cmdRegister:", err)
		return
	}
	do.Color = 0x3f82ef
	do.Description = "You have successfully registered."
	return
}

// The roles command receives a Discord session pointer, a channel, a user and an arguments slice of strings.
// It then shows a list of added and available roles or allows the user to add or remove roles on the server.
func cmdRoles(dg *discordgo.Session, channel string, user string, args []string) (do *DiscordOutput) {
	do = NewDiscordOutput(dg, 0xb40000, "ROLES", "")
	users, err := readCSV(usersFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		log.Println("cmdRoles:", err)
		return
	}
	for _, u := range users {
		if strings.EqualFold(u[0], user) {
			if strings.Contains(strings.ToLower(u[2]), "embeds") {
				do.Embeds = true
			}
		}
	}
	roles, err := readCSV(rolesFile)
	if err != nil {
		do.Description = ":warning: Error getting roles."
		log.Println("cmdRoles:", err)
		return
	}
	guildRoles, err := dg.GuildRoles(guild)
	if err != nil {
		do.Description = ":warning: Error getting guild roles."
		log.Println("cmdRoles:", err)
		return
	}
	member, err := dg.GuildMember(guild, user)
	if err != nil {
		do.Description = ":warning: Error getting guild member."
		log.Println("cmdRoles:", err)
		return
	}
	// We have two possibilities, either the user passes no arguments or else, we consider only the first one.
	// The first branch doesn't have early exits, so we use the return at the end of the function to return the output.
	// The second branch though has several possible early exists, so we use local returns to return the output.
	// Except at the very end, where nothing else will run, so we can use the return at the end of the function.
	if len(args) == 0 {
		do.Color = 0x3f82ef
		do.Description += "**Current roles you are added to:**\n\n"
		for _, v := range member.Roles {
			for _, w := range guildRoles {
				if strings.EqualFold(v, w.ID) {
					do.Description += w.Name + "\n"
				}
			}
		}
		do.Description += "\n**Available roles that you can add/remove:**\n\n"
		for _, v := range roles {
			do.Description += v[0] + "\n"
		}
		do.Description += "\n**To add/remove roles call this command with a role name.**\n\nExample: !roles space_notifications"
	} else {
		valid := false
		for _, v := range roles {
			if strings.EqualFold(args[0], v[0]) {
				valid = true
			}
		}
		if !valid {
			do.Description = ":warning: This is not one of the available roles."
			return
		}
		for _, v := range guildRoles {
			if strings.EqualFold(args[0], v.Name) {
				for _, w := range member.Roles {
					if strings.EqualFold(v.ID, w) {
						// The user is already assigned to this role.
						err := dg.GuildMemberRoleRemove(guild, user, v.ID)
						if err != nil {
							do.Description = ":warning: Error removing role."
							log.Println("cmdRoles:", err)
							return
						}
						do.Description = fmt.Sprintf("You were successfully removed from the %s role.", v.Name)
						return
					}
				}
				// The user is not assigned to this role yet.
				err := dg.GuildMemberRoleAdd(guild, user, v.ID)
				if err != nil {
					do.Description = ":warning: Error adding role."
					log.Println("cmdRoles:", err)
					return
				}
				do.Description = fmt.Sprintf("You were successfully added to the %s role.", v.Name)
				return
			}
		}
		// We only get here, if one of the allowed roles from the CSV file was not found on the server.
		do.Description = fmt.Sprintf("The %s role doesn't exist on the server.", args[0])
	}
	return
}

// The stats command receives a Discord session pointer, a channel, and a user.
// It then reads some general user stats periodically stored and displays them.
func cmdStats(dg *discordgo.Session, channel string, user string) {
	do := NewDiscordOutput(dg, 0xb40000, "STATS", "")
	stats, err := readCSV(statsFile)
	if err != nil {
		log.Println("tskStats:", err)
		return
	}
	// Go through all the stats lines and append each one as a chart value.
	var values []chart.Value
	for _, v := range stats {
		valueFloat, _ := strconv.ParseFloat(v[1], 64)
		member, _ := dg.GuildMember(guild, v[0])
		label := ""
		if valueFloat > 40 {
			label = member.User.Username + " - " + v[1]
		}
		values = append(values, chart.Value{Value: valueFloat, Label: label})
	}
	pie := chart.PieChart{
		Title:      "Total Messages",
		TitleStyle: chart.StyleShow(),
		Width:      700,
		Height:     800,
		Values:     values,
	}
	f, _ := os.Create("messages.png")
	pie.Render(chart.PNG, f)
	f.Close()
	f, _ = os.Open("messages.png")
	do.File(channel, "messages.png", f, "**STATS**")
	f.Close()
}

// The weather command receives a Discord session pointer, a channel, a user and an arguments slice of strings.
// It then shows the current weather for a given location on the channel using the OpenWeatherMap API.
func cmdWeather(dg *discordgo.Session, channel string, user string, args []string) (do *DiscordOutput) {
	do = NewDiscordOutput(dg, 0xb40000, "WEATHER", "")
	users, err := readCSV(usersFile)
	if err != nil {
		do.Description = ":warning: Error getting users."
		log.Println("cmdWeather:", err)
		return
	}
	for _, u := range users {
		if strings.EqualFold(u[0], user) {
			if strings.Contains(strings.ToLower(u[2]), "embeds") {
				do.Embeds = true
			}
		}
	}
	weather, err := readCSV(weatherFile)
	if err != nil {
		do.Description = ":warning: Error getting weather settings."
		log.Println("cmdWeather:", err)
		return
	}
	location := ""
	tempUnits := "C"
	windUnits := "m/s"
	// Neither a location nor temperature unit were provided as an argument to the command.
	// So we must get the location and temperature unit for the user from the weather file.
	// If a user in the weather file matches user, we get its location and temperature unit.
	if len(args) == 0 {
		for _, v := range weather {
			if strings.EqualFold(v[0], user) {
				tempUnits = strings.ToUpper(v[1])
				location = v[2]
			}
		}
		// A temperature unit was provided as an argument to the command, we must update the setting.
		// However, we must first check if the user already has a location set on the weather file.
		// If so, we update the user units, otherwise we ask him to get the weather for a location.
		// This is so that the user gets registered on the weather file before we can set a location.
	} else if len(args) == 1 && (strings.ToLower(args[0]) == "c" || strings.ToLower(args[0]) == "f") {
		var unitsUpdated bool
		for i, v := range weather {
			// User with a location on the weather database.
			if strings.EqualFold(v[0], user) {
				unitsUpdated = true
				weather[i][1] = strings.ToLower(args[0])
			}
		}
		if !unitsUpdated {
			do.Description = ":warning: Get the weather for some location before setting the units."
			return
		}
		err = writeCSV(weatherFile, weather)
		if err != nil {
			do.Description = ":warning: Error storing weather units."
			log.Println("cmdWeather:", err)
			return
		}
		do.Description = "Temperature units updated."
		return
		// If we reach this point, a location was provided as an argument to the command.
		// If the user already exists, we update his location, otherwise we register him.
	} else {
		var newUser bool = true
		location = strings.Join(args, " ")
		for i, v := range weather {
			// User with a location on the weather database.
			if strings.EqualFold(v[0], user) {
				newUser = false
				weather[i][2] = location
			}
		}
		if newUser {
			// User without a location on the weather database.
			weather = append(weather, []string{user, "c", location})
		}
		err = writeCSV(weatherFile, weather)
		if err != nil {
			do.Description = ":warning: Error storing weather location."
			log.Println("cmdWeather:", err)
			return
		}
	}
	if location == "" {
		do.Description = ":warning: Please provide a location as argument."
		return
	}
	if tempUnits == "F" {
		windUnits = "mph"
	}
	// Finally we get the current weather at a location using the temperature units.
	// Then we display a nicely formatted and compact weather string on the channel.
	w, err := owm.NewCurrent(tempUnits, "en", owmAPIKey)
	if err != nil {
		do.Description = ":warning: Error fetching weather."
		log.Println("cmdWeather:", err)
		return
	}
	err = w.CurrentByName(location)
	if err != nil {
		do.Description = ":warning: Could not fetch weather for that location."
		log.Println("cmdWeather:", err)
		return
	}
	icon := ""
	switch {
	case strings.Contains(strings.ToLower(w.Weather[0].Description), "clear sky"):
		icon = ":sunny:"
	case strings.Contains(strings.ToLower(w.Weather[0].Description), "few clouds"):
		icon = ":white_sun_small_cloud:"
	case strings.Contains(strings.ToLower(w.Weather[0].Description), "broken clouds"):
		icon = ":white_sun_cloud:"
	case strings.Contains(strings.ToLower(w.Weather[0].Description), "scattered clouds"):
		icon = ":white_sun_cloud:"
	case strings.Contains(strings.ToLower(w.Weather[0].Description), "overcast clouds"):
		icon = ":cloud:"
	case strings.Contains(strings.ToLower(w.Weather[0].Description), "rain"):
		icon = ":cloud_rain:"
	case strings.Contains(strings.ToLower(w.Weather[0].Description), "fog"):
		icon = ":fog:"
	}
	icon += " " + w.Weather[0].Description
	do.Color = 0x3f82ef
	do.Description =
		fmt.Sprintf("**%s**\n\n%s\n\n:thermometer: %0.1f%s\n:droplet: %d%%\n:arrow_down: %0.1fhPa\n:triangular_flag_on_post: %0.1f%s",
			w.Name,
			icon,
			w.Main.Temp,
			tempUnits,
			w.Main.Humidity,
			w.Main.Pressure,
			w.Wind.Speed,
			windUnits)
	return
}

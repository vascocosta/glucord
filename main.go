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
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	prefix       = "!" // Prefix which is used by the user to issue commands.
	token        = ""  // Token used to authenticate the bot with Discord.
	guild        = ""  // Guild ID.
	feedInterval = 300 // Feed poll interval in seconds.
)

const (
	answersFile   = "answers.csv" // Full path to the answers file.
	configFile    = "config.csv"  // Full path to the config file.
	eventsFile    = "events.csv"  // Full path to the events file.
	feedsFile     = "feeds.csv"   // Full path to the feeds file.
	pluginsFolder = "./plugins/"  // Full path to the plugins folder.
	quotesFile    = "quotes.csv"  // Full path to the quotes file.
	usageFile     = "usage.csv"   // Full path to the usage file.
	usersFile     = "users.csv"   // Full path to the users file.
	hns           = 3600000000000 // Number of nanoseconds in one hour.
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	m.Content = strings.Trim(m.Content, " ")
	command, err := parseCommand(m.Content, m.Author.ID, m.ChannelID)
	if err != nil {
		return
	} else {
		var do *DiscordOutput
		switch strings.ToLower(command.Name) {
		case "a", "ask":
			do = cmdAsk(s, command.Channel, command.User, command.Args)
		case "h", "help", "commands":
			do = cmdHelp(s, command.Channel, command.User, strings.Join(command.Args, ""))
		case "n", "next":
			do = cmdNext(s, command.Channel, command.User, strings.Join(command.Args, " "))
		case "p", "ping":
			do = cmdPing(s, command.Channel, command.User, command.Args)
		case "q", "quote":
			cmdQuote(s, command.Channel, command.User, command.Args)
		default:
			go cmdPlugin(strings.ToLower(command.Name), s, command.Channel, command.User, command.Args)
			return
		}
		if do != nil {
			if do.Embeds {
				s.ChannelMessageSendEmbed(command.Channel, do.Embed())
			} else {
				s.ChannelMessageSend(command.Channel, do.Text())
			}
		}
	}

}

func main() {
	config, err := readCSV(configFile)
	if err != nil {
		log.Println("main:", err)
		return
	}
	prefix = config[0][0]
	token = config[0][1]
	guild = config[0][2]
	feedInterval, _ = strconv.Atoi(config[0][3])
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("main:", err)
		return
	}
	dg.AddHandler(messageCreate)
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	err = dg.Open()
	if err != nil {
		fmt.Println("main:", err)
		return
	}
	go tskFeeds(dg)
	go tskEvents(dg)
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, guild, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	for _, v := range registeredCommands {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, guild, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
	dg.Close()
}

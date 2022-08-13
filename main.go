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
	feedInterval = 300 // Feed poll interval in seconds.
)

const (
	answersFile   = "answers.csv" // Full path to the answers file.
	configFile    = "config.csv"  // Full path to the config file.
	eventsFile    = "events.csv"  // Full path to the events file.
	feedsFile     = "feeds.csv"   // Full path to the feeds file.
	pluginsFolder = "./plugins/"  // Full path to the plugins folder.
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
		switch strings.ToLower(command.Name) {
		case "a", "ask":
			cmdAsk(s, command.Channel, command.User, command.Args)
		case "h", "help", "commands":
			cmdHelp(s, command.Channel, command.User, strings.Join(command.Args, ""))
		case "n", "next":
			cmdNext(s, command.Channel, command.User, strings.Join(command.Args, " "))
		default:
			go cmdPlugin(strings.ToLower(command.Name), s, command.Channel, command.User, command.Args)
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
	feedInterval, _ = strconv.Atoi(config[0][2])
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("main:", err)
		return
	}
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	err = dg.Open()
	if err != nil {
		fmt.Println("main:", err)
		return
	}
	fmt.Println("The bot is now running. Press CTRL-C to exit.")
	go tskFeeds(dg)
	go tskEvents(dg)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

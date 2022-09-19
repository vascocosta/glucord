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
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	prefix       = "!" // Prefix which is used by the user to issue commands.
	token        = ""  // Token used to authenticate the bot with Discord.
	guild        = ""  // Guild ID.
	feedInterval = 300 // Feed poll interval in seconds.
	owmAPIKey    = ""  // OWM API key.
)

const (
	aliasFile     = "alias.csv"    // Full path to the alias file.
	answersFile   = "answers.csv"  // Full path to the answers file.
	betFile       = "bet.csv"      // Full path to the bet file.
	betsFile      = "bets.csv"     // Full path to the bets file.
	disabledFile  = "disabled.csv" // Full path to the disabled file.
	configFile    = "config.csv"   // Full path to the config file.
	driversFile   = "drivers.csv"  // Full path to the drivers file.
	eventsFile    = "events.csv"   // Full path to the events file.
	feedsFile     = "feeds.csv"    // Full path to the feeds file.
	inputFile     = "input.txt"    // Full path to the input file.
	pluginsFolder = "./plugins/"   // Full path to the plugins folder.
	quotesFile    = "quotes.csv"   // Full path to the quotes file.
	resultsFile   = "results.csv"  // Full path to the results file.
	rolesFile     = "roles.csv"    // Full path to the roles file.
	statsFile     = "stats.csv"    // Full path to the stats file.
	usageFile     = "usage.csv"    // Full path to the usage file.
	usersFile     = "users.csv"    // Full path to the users file.
	weatherFile   = "weather.csv"  // Full path to the weather file.
	hns           = 3600000000000  // Number of nanoseconds in one hour.
)

// Message callback function that receives a Discord session pointer and a message pointer.
// If there is a command at the beginning of a message, it runs the corresponding function.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	m.Content = strings.Trim(m.Content, " ")
	command, err := parseCommand(m.Content, m.Author.ID, m.ChannelID)
	if err != nil {
		return
	} else {
		// Pick the corresponding function for each supported command and store its output.
		// If the command is not built-in, run it as a plugin inside a dedicated goroutine.
		disabled, err := readCSV(disabledFile)
		if err != nil {
			log.Println("main:", err)
		}
		for _, v := range disabled {
			if strings.EqualFold(v[0], command.Name) {
				s.ChannelMessageSend(command.Channel, ":warning: Unkown command or plugin.")
				return
			}
		}
		var do *DiscordOutput
		switch strings.ToLower(command.Name) {
		case "a", "ask":
			do = cmdAsk(s, command.Channel, command.User, command.Args)
		case "b", "bet":
			do = cmdBet(s, command.Channel, command.User, command.Args)
		case "h", "help", "commands":
			do = cmdHelp(s, command.Channel, command.User, strings.Join(command.Args, ""))
		case "n", "next":
			do = cmdNext(s, command.Channel, command.User, strings.Join(command.Args, " "))
		case "p", "ping":
			do = cmdPing(s, command.Channel, command.User, command.Args)
		case "pb", "processbets":
			do = cmdProcessBets(s, command.Channel, command.User)
		case "q", "quote":
			do = cmdQuote(s, command.Channel, command.User, command.Args)
		case "r", "register":
			do = cmdRegister(s, command.Channel, command.User)
		case "ro", "roles":
			do = cmdRoles(s, command.Channel, command.User, command.Args)
		case "s", "stats":
			cmdStats(s, command.Channel, command.User)
		case "w", "weather":
			do = cmdWeather(s, command.Channel, command.User, command.Args)
		default:
			finishedCh := make(chan bool)
			go func() {
				select {
				case <-finishedCh:
					// We received true from the finished channel, so we do nothing.
				case <-time.After(3 * time.Second):
					// We didn't receive from the finished channel within 3 seconds.
					s.ChannelMessageSend(command.Channel, ":warning: Command is taking long to run... Please wait.")
				}
			}()
			go cmdPlugin(strings.ToLower(command.Name), s, command.Channel, command.User, command.Args, finishedCh)
		}
		// If the pointer to DiscordOutput isn't nil (built-in command) send the output here.
		// This is to prevent access to methods on a nil pointer (cmdPlugin does not set do).
		if do != nil {
			if do.Embeds {
				s.ChannelMessageSendEmbed(command.Channel, do.Embed())
			} else {
				s.ChannelMessageSend(command.Channel, do.Text())
			}
		}
	}

}

// The main function initialises some variables from a configuration file, then sets up the bot and connects to Discord.
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
	owmAPIKey = config[0][4]
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("main:", err)
		return
	}
	// Add callback function to handle messages and fire up the appropriate regular command function.
	dg.AddHandler(messageCreate)
	// Add callback function to handle interactions and fire up the appropriate slash command function.
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Slash command names are mapped to corresponding handler functions on the commandHandlers variable.
		// If the name of the slash command is a valid key of commandHandlers, execute the handler function.
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	err = dg.Open()
	if err != nil {
		log.Println("main:", err)
		return
	}
	// Launch some background tasks that run concurrently with the main goroutine of the bot.
	// These functions need to keep running in the background the whole time to perform work.
	// While bot commands are user triggered and short lived these tasks happen periodically.
	go tskEvents(dg)
	go tskFeeds(dg)
	go tskStats(dg)
	go tskWrite(dg)
	// Keep a record of all the slash commands defined in the commands variable using a slice.
	// Register a slash command on Discord for every command defined in the commands variable.
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, guild, v)
		if err != nil {
			log.Println("main:", err)
		}
		registeredCommands[i] = cmd
	}
	// The work of the main goroutine of the bot, which is to set it up, is done by this point.
	// However, we need to prevent it from finishing and kill all other goroutines prematurely.
	// For this, we redirect any termination signals to a channel using the os.Notify function.
	// The main goroutine will block reading the channel and exit only when one signal arrives.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	// Unregister a slash command on Discord for every command defined in the commands variable.
	// This runs after a termination signal is caught to delete any slash commands from Discord.
	for _, v := range registeredCommands {
		err := dg.ApplicationCommandDelete(dg.State.User.ID, guild, v.ID)
		if err != nil {
			log.Println("main:", err)
		}
	}
	dg.Close()
}

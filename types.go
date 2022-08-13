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
	"github.com/bwmarrin/discordgo"
)

// Type that represents a Discord command issued by the user.
type Command struct {
	Name    string
	Args    []string
	User    string
	Channel string
}

type DiscordOutput struct {
	Title string
	Description string
	Color int
	Session *discordgo.Session
	Embeds bool
}

func NewDiscordOutput(title string, description string, color int, s *discordgo.Session, embeds bool) *DiscordOutput {
	return &DiscordOutput{title, description, color, s, embeds}
}

func (do *DiscordOutput) Send(channel string) {
	if do.Embeds {
		output := &discordgo.MessageEmbed{}
		output.Title = do.Title
		output.Description = do.Description
		output.Color = do.Color
		do.Session.ChannelMessageSendEmbed(channel, output)
	} else {
		do.Session.ChannelMessageSend(channel, do.Description)
	}
}
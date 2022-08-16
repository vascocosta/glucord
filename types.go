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
	Title       string
	Description string
	Color       int
	Session     *discordgo.Session
	Embeds      bool
	Fields      *[]map[string]string
	Image       *string
}

func NewDiscordOutput(title string, description string, color int, s *discordgo.Session, embeds bool, fields *[]map[string]string, image *string) *DiscordOutput {
	return &DiscordOutput{title, description, color, s, embeds, fields, image}
}

func (do *DiscordOutput) Send(channel string) {
	if do.Embeds {
		output := &discordgo.MessageEmbed{}
		output.Title = do.Title
		if do.Fields == nil {
			output.Description = do.Description
		}
		output.Color = do.Color
		if do.Fields != nil {
			for _, v := range *do.Fields {
				field := &discordgo.MessageEmbedField{}
				field.Name = v["Name"]
				field.Value = v["Value"]
				output.Fields = append(output.Fields, field)
			}
		}
		if do.Image != nil {
			embedImage := &discordgo.MessageEmbedImage{}
			embedImage.URL = *do.Image
			output.Image = embedImage
		}
		do.Session.ChannelMessageSendEmbed(channel, output)
	} else {
		if do.Fields != nil {
			for _, v := range *do.Fields {
				do.Description += fmt.Sprintf("**%s**\n%s\n", v["Name"], v["Value"])
			}
		}
		if do.Image != nil {
			do.Description += *do.Image
		}
		do.Session.ChannelMessageSend(channel, do.Description)
	}
}

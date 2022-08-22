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
	"encoding/csv"
	"errors"
	"os"
	"strings"
)

// Small utility function that returns weather a slice of strings contains a given string.
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// Small utility function that returns whether a user exists or not.
func isUser(user string, users [][]string) bool {
	for _, v := range users {
		if strings.EqualFold(user, v[0]) {
			return true
		}
	}
	return false
}

// Small utility function that reads a CSV file and returns the data as slice of slice of strings.
func readCSV(path string) (data [][]string, err error) {
	f, err := os.Open(path)
	if err != nil {
		err = errors.New("Error opening CSV file: " + path + ".")
		return
	}
	defer f.Close()
	r := csv.NewReader(f)
	data, err = r.ReadAll()
	if err != nil {
		err = errors.New("Error reading data from: " + path + ".")
		return
	}
	return
}

// Small utility function that writes a slice of slice of strings to a CSV file.
func writeCSV(path string, data [][]string) (err error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		err = errors.New("Error opening CSV file: " + path + ".")
		return
	}
	defer f.Close()
	w := csv.NewWriter(f)
	err = w.WriteAll(data)
	if err != nil {
		err = errors.New("Error writing data to: " + path + ".")
		return
	}
	return
}

// Small utility function that takes a message string and breaks it down into a Command.
func parseCommand(message string, user string, channel string) (command Command, err error) {
	if len(message) > 1 && strings.HasPrefix(message, prefix) {
		split := strings.Split(message, " ")
		command.Name = split[0][1:]
		command.Args = split[1:]
		command.User = user
		command.Channel = channel
		return
	} else {
		err = errors.New("invalid command")
		return
	}
}

// Small utility function that checks if a file exists and is not a directory.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

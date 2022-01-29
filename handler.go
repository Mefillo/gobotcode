package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ansel1/merry"
)

type Data struct {
	Response string
}

var COMMANDS = map[string]bool{"fl": true, "fa": true, "fd": true}

const CANCEL = "c"

func processRequest(update Update) (data Data, err error) {
	fmt.Printf("UPDATE")
	fmt.Printf("%+v", update)

	// Sanitize input
	var sanitizedSeed = sanitize(update.Message.Text)
	var response string

	// Get current record
	item, err := Get(update.Message.From.Username)
	if err != nil {
		fmt.Printf("\n !!! Got error getting: %+v\n", err)
	}

	// Check if there is conv in progress
	if item.Status != "" {
		fmt.Printf("Current status not nil: %s", item.Status)
		switch item.Status {
		case "fa":
			if !stringInSlice(sanitizedSeed, item.Films) && sanitizedSeed != CANCEL {
				item.Films = append(item.Films, sanitizedSeed)
			}
			item.Status = ""
			err = Save(item)
			if err != nil {
				fmt.Printf("Got error saving data to DB: %+v", err)
				return
			}
			response = "ok"
		case "fd":
			item.Status = ""
			indexToDelete, e := strconv.Atoi(sanitizedSeed)
			if e != nil {
				response = "meh"
			} else {
				if indexToDelete >= len(item.Films) && indexToDelete < 0 {
					response = "hold your hourses pal"
				} else {
					response = item.Films[indexToDelete]
				}
			}
			err = Save(item)
			if err != nil {
				fmt.Printf("Got error saving data to DB: %+v", err)
				return
			}
		}
	} else {
		// Check for commands actions
		com := strings.ToLower(sanitizedSeed)
		if COMMANDS[com] {
			return commands_handler(item, com)
		}
		response = "?"
	}

	d := Data{Response: response}
	return d, err

}

// sanitize remove clutter like /start /punch or the bot name from the string s passed as input
func sanitize(s string) string {
	if len(s) >= lenStartCommand {
		if s[:lenStartCommand] == startCommand {
			s = s[lenStartCommand:]
		}
	}

	if len(s) >= lenPunchCommand {
		if s[:lenPunchCommand] == punchCommand {
			s = s[lenPunchCommand:]
		}
	}
	if len(s) >= lenBotTag {
		if s[:lenBotTag] == botTag {
			s = s[lenBotTag:]
		}
	}

	return s
}

func commands_handler(item Item, status string) (data Data, err error) {
	switch status {
	case "fl":
		data.Response = strings.Join(item.Films, ", ")
		return
	case "fa":
		// save new status
		item.Status = "fa"
		err = Save(item)
		if err != nil {
			return
		}
		// ask for new movie for the list
		data.Response = "Add new movie title:"
		return
	case "fd":
		// save new status
		item.Status = "fd"
		err = Save(item)
		if err != nil {
			return
		}
		// ask for new movie for the list
		data.Response = convertToIndexed(item.Films)
		return
	}

	err = merry.New(fmt.Sprintf("No such action for one level handler: %+v", status))
	return
}

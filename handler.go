package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ansel1/merry"
)

type Data struct {
	Response string
}

var ONE_LEVEL = map[string]bool{"fl": true}

func processRequest(update Update) (data Data, err error) {
	fmt.Printf("UPDATE")
	fmt.Printf("%+v", update)

	// Sanitize input
	var sanitizedSeed = sanitize(update.Message.Text)

	// Get current record
	item, err := Get(update.Message.From.Username)
	if err != nil {
		fmt.Printf("\n !!! Got error getting: %+v\n", err)
	}

	// Check for one level actions
	if ONE_LEVEL[item.Status] {
		return one_level_handler(item, item.Status)
	}

	// Change item
	item.Films = append(item.Films, sanitizedSeed)

	// Save response to DB
	stringItem := fmt.Sprintf(`{
		"id": "%s",
		"film": "%s"
	}`, update.Message.From.Username, update.Message.Text)

	err = json.Unmarshal([]byte(stringItem), &item)
	if err != nil {
		fmt.Printf("Got error unmarshaling data: %+v", err)
		return
	}

	err = Save(item)
	var response string
	if err != nil {
		fmt.Printf("Got error saving data to DB: %+v", err)
		return
	} else {
		response = fmt.Sprintf("New film list: %+v", item.Films)
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

func one_level_handler(item Item, status string) (data Data, err error) {
	switch status {
	case "fl":
		data.Response = strings.Join(item.Films, ", ")
		return
	}
	err = merry.New(fmt.Sprintf("No such action for one level handler: %+v", status))
	return
}

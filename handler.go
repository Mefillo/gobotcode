package main

import (
	"encoding/json"
	"fmt"
)

type Data struct {
	Response string
}

func processRequest(update Update) (err error, data Data) {
	fmt.Printf("UPDATE")
	fmt.Printf("%+v", update)

	// Get current record
	var item Item
	item, err = Get(update.Message.From.Username)
	if err != nil {
		fmt.Printf("\n !!! Got error getting: %+v\n", err)
	}

	// Change item
	// Sanitize input
	var sanitizedSeed = sanitize(update.Message.Text)
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
	return err, d

}

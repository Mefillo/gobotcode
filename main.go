// Package handler contains an HTTP Cloud Function to handle update from Telegram whenever a users interacts with the
// bot.

// https://api.telegram.org/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/getMe
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Define a few constants and variable to handle different commands
const punchCommand string = "/punch"

var lenPunchCommand int = len(punchCommand)

const startCommand string = "/start"

var lenStartCommand int = len(startCommand)

const botTag string = "@RapGeniusBot"

var lenBotTag int = len(botTag)

// Pass token and sensible APIs through environment variables
const telegramApiBaseUrl string = "https://api.telegram.org/bot"
const telegramApiSendMessage string = "/sendMessage"

var telegramApi string = telegramApiBaseUrl + os.Getenv("BOT_TOKEN") + telegramApiSendMessage

// Update is a Telegram object that we receive every time an user interacts with the bot.
type Update struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

type Response struct {
	Message string `json: "Answer:"`
}

// HandleTelegramWebHook sends a message back to the chat with a punchline starting by the message provided by the user.
func HandleTelegramWebHook(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse and validate request
	update, err := parseTelegramRequest(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(`{"ok":"nope"}`),
		}, nil
	}
	if update.Message.From.Username == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       `{"success": true}`,
		}, nil
	}
	data, err := processRequest(*update)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusExpectationFailed,
			Body:       fmt.Sprintf("error from process: %+v", err),
		}, nil
	}
	// Send response back to Telegram
	var telegramResponseBody, errTelegram = sendTextToTelegramChat(update.Message.Chat.ID, data.Response)
	fmt.Printf("SENDING RESUTLS: %+v !!! %+v", telegramResponseBody, errTelegram)
	if errTelegram != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusExpectationFailed,
			Body:       string("error from telegram"),
		}, nil
	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       `{"success": true}`,
		}, nil
	}
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(r events.APIGatewayProxyRequest) (*Update, error) {
	// var wht interface{}
	// json.Unmarshal([]byte(r.Body), &wht)
	// fmt.Printf("!!!WHATVR: %+v\n", wht)
	var update Update
	err := json.Unmarshal([]byte(r.Body), &update)
	if err != nil {
		fmt.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}
	if update.UpdateID == 0 {
		fmt.Printf("invalid update id, got update id = 0")
		return nil, errors.New("invalid update id of 0")
	}
	return &update, nil
}

func main() {
	lambda.Start(HandleTelegramWebHook)
}

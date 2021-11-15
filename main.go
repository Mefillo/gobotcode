// Package handler contains an HTTP Cloud Function to handle update from Telegram whenever a users interacts with the
// bot.

// https://api.telegram.org/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/getMe
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

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

type Item struct {
	ID   string `json:"id"`
	Films []string `json:"films"`
}

type Response struct {
	Message string `json: "Answer:"`
}

// HandleTelegramWebHook sends a message back to the chat with a punchline starting by the message provided by the user.
func HandleTelegramWebHook(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	update, err := parseTelegramRequest(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(`{"ok":"nope"}`),
		}, nil
	}
	fmt.Printf("UPDATE")
	fmt.Printf("%+v", update)
	// Sanitize input
	var sanitizedSeed = sanitize(update.Message.Text)

	// Get response
	var response = sanitizedSeed + update.Message.From.Username

	// Get current record
	var item Item
	item, err = Get(update.Message.From.Username)
	if err != nil {
		fmt.Printf("\n !!! Got error getting: %+v\n", err)
	}

	// Change item
	item.Films = append(item.Films, update.Message.Text)

	// Save response to DB
	stringItem := fmt.Sprintf(`{
		"id": "%s",
		"film": "%s"
	}`, update.Message.From.Username, update.Message.Text)
	err = json.Unmarshal([]byte(stringItem), &item)
	err = Save(item)
	if err != nil {
		fmt.Printf("Got error saving data to DB: %+v", err)
	}

	// Send response back to Telegram
	var telegramResponseBody, errTelegram = sendTextToTelegramChat(update.Message.Chat.ID, response)
	fmt.Printf("SENDING RESUTLS: %+v !!! %+v", telegramResponseBody, errTelegram)
	if errTelegram != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusExpectationFailed,
			Body:       string("error from telegram"),
		}, nil
	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string("ok"),
		}, nil
	}
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(r events.APIGatewayProxyRequest) (*Update, error) {
	var update Update
	err := json.Unmarshal([]byte(r.Body), &update)
	if err != nil {
		fmt.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}
	if update.UpdateID == 0 {
		fmt.Printf("invalid update id, got update id = 0")
		return nil, errors.New("invalid update id of 0 indicates failure to parse incoming update")
	}
	return &update, nil
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

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendTextToTelegramChat(chatId int, text string) (string, error) {

	fmt.Printf("Sending %s to chat_id: %d", text, chatId)
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

	if err != nil {
		fmt.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		fmt.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	fmt.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}

func main() {
	lambda.Start(HandleTelegramWebHook)
}

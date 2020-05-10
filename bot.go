package main

import (
	"bytes"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Quiz struct
type Quiz struct {
	ChatID          int64    `json:"chat_id"`
	Question        string   `json:"question"`
	Options         []string `json:"options"`
	Type            string   `json:"type"`
	CorrectOptionID int      `json:"correct_option_id"`
}

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/next"),
	),
)

func telegramBot(token string, hook string, cert string, key string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	data, err := GetData()
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 60
	// updates, err := bot.GetUpdatesChan(u)

	_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert(fmt.Sprintf("https://%s/%s", hook, token), cert))
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS("0.0.0.0:8443", cert, key, nil)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		if update.CallbackQuery != nil {
			fmt.Print(update)

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			opt := []string{"Test1", "TEst2", "Test3"}
			quiz := Quiz{ChatID: update.Message.Chat.ID, Question: "Choose right answer:", Options: opt, Type: "quiz", CorrectOptionID: 0}
			switch update.Message.Command() {
			case "start":
				msg.Text = "Type /next."
			case "help":
				msg.Text = "Type /next."
			case "next":
				rand.Seed(time.Now().UTC().UnixNano())
				rand := rand.Intn(len(data))
				res := data[rand]
				msg.Text = res.Origin + " - " + res.Translate
				msg.ReplyMarkup = numericKeyboard
				url := "https://api.telegram.org/bot1209313230:AAFK2qDwS7SKnrWDXFxdmZVQYuw6CYNkoMg/sendPoll"
				var jsonStr = []byte(fmt.Sprintf("%v", quiz))
				log.Printf("poll: %s", fmt.Sprintf("%v", quiz))
				req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
				if err != nil {
					panic(err)
				}
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()
			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
		}
	}
}

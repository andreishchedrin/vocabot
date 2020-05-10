package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math/rand"
	"time"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/next"),
	),
)

func telegramBot(token string) {
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

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

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
			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
		}
	}
}

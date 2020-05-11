package main

import (
	"bytes"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math/rand"
	"net/http"
	"time"
	"encoding/json"
	"strings"
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

func getRand(length int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	index := rand.Intn(length)
	return index
}

func getRandSlice(slice *[]string) {
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Shuffle(len(*slice), func(i, j int) {
		(*slice)[i], (*slice)[j] = (*slice)[j], (*slice)[i]
	})
}

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

		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {		

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			
			switch update.Message.Command() {
			case "start":
				msg.Text = "Type /next."
			case "help":
				msg.Text = "Type /next."
			case "next":
				
				// Choose the word which use as rigth answer
				rand := getRand(len(data))				
				rightAnswer := data[rand]
				var rightAnswerIndex int

				// Choose fake answers
				list := []string{
					data[getRand(len(data))].Translate,
					data[getRand(len(data))].Translate,
					data[getRand(len(data))].Translate,
					rightAnswer.Translate,
				}
				getRandSlice(&list)				
				for i, listItem := range list {
					if listItem == rightAnswer.Translate {
						rightAnswerIndex = i
					}
				}

				// Prepare sendPoll API request 
				str := strings.ToUpper(rightAnswer.Origin) + ":"				
				quiz := &Quiz{ChatID: update.Message.Chat.ID, Question: str, Options: list, Type: "quiz", CorrectOptionID: rightAnswerIndex}
				url := "https://api.telegram.org/bot1209313230:AAFK2qDwS7SKnrWDXFxdmZVQYuw6CYNkoMg/sendPoll"
				buf := new(bytes.Buffer)
				json.NewEncoder(buf).Encode(quiz)	
				// Send poll request			
				req, err := http.NewRequest("POST", url, buf)
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

				// Log for debug
				// log.Printf("poll: %s", fmt.Sprintf("%v", quiz))
				
				// msg.Text = res.Origin + " - " + res.Translate
				// msg.ReplyMarkup = numericKeyboard
				msg.Text = "Tap /next to continue.."
			default:
				msg.Text = "I don't know that command"
			}

			msg.ReplyMarkup = numericKeyboard
			bot.Send(msg)
				
		}
	}
}

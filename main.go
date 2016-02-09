package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
)

var tokenChatMap map[string]int
var bot *tgbotapi.BotAPI

func handler(w http.ResponseWriter, r *http.Request) {
	token := fmt.Sprint(r.URL.Path[1:])
	chatID := tokenChatMap[token]

	text := r.FormValue("text")
	if text == "" {
		text = "<empty>"
	}

	bot.Send(tgbotapi.NewMessage(chatID, text))
	fmt.Fprintf(w, "Token: %s, Chat: %d!", token, chatID)
}

func main() {
	var err error
	bot, err = tgbotapi.NewBotAPI("180160051:AAECYIkVGlurwZwO4Y7KgPN8T9jRNHnxYSs")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
		return
	}

	tokenChatMap = make(map[string]int)

	go func() {
		for update := range updates {
			// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if update.Message.Text == "/start" {
				log.Printf("%d - %s", update.Message.Chat.ID, update.Message.From.UserName)
				token := fmt.Sprint(uuid.NewV4())
				tokenChatMap[token] = update.Message.Chat.ID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Hello %s, your token is %s", update.Message.From.FirstName, token))
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Yo!!!")
				bot.Send(msg)
			}
		}
	}()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":9090", nil)
}

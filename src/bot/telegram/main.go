package main

import (
	"log"
	"os"

	TelegramApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	var botError error
	token := os.Getenv("TOKEN_TELEGRAM")
	bot, botError := TelegramApi.NewBotAPI(token)

	if botError != nil {
		log.Panic(botError)
	}

	botUpdate := TelegramApi.NewUpdate(0)
	botUpdate.Timeout = 60

	updates := bot.GetUpdatesChan(botUpdate)
	var message string

	for update := range updates {
		if update.ChannelPost != nil {
			log.Printf("[Channel - %s] %s", update.ChannelPost.SenderChat.Title, update.ChannelPost.Text)
			message = update.ChannelPost.Text
		} else if update.Message != nil {
			log.Printf("[User - %s] %s", update.Message.From.UserName, update.Message.Text)
			message = update.Message.Text
		}

		// TODO: send message to API
		log.Printf("Message: %s", message)
	}
}

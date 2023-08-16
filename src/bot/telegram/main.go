package main

import (
	"log"

	TelegramApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	var botError error
	bot, botError := TelegramApi.NewBotAPI("5737642054:AAFH6ULajxnK59ySDCCRzKLxyKgkXfPImWM")

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

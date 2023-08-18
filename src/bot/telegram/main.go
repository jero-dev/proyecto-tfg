package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := os.Getenv("TOKEN_TELEGRAM")
	bot, botError := telegram.NewBotAPI(token)

	if botError != nil {
		log.Panic(botError)
	}

	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	http.HandleFunc("/api/setup", func(w http.ResponseWriter, r *http.Request) {
		setUpWebhook(w, r, bot)
	})
	http.HandleFunc("/api/handleUpdate", handleTelegramUpdate)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func setUpWebhook(w http.ResponseWriter, r *http.Request, bot *telegram.BotAPI) {

	webHook, _ := telegram.NewWebhook("https://" + r.Host + "/api/handleUpdate")

	_, webHookError := bot.Request(webHook)

	if webHookError != nil {
		log.Fatal(webHookError)
	}

	_, printError := fmt.Fprint(w, "Se ha configurado el webhook correctamente")
	if printError != nil {
		return
	}
}

func handleTelegramUpdate(w http.ResponseWriter, r *http.Request) {

	update := &telegram.Update{}
	if decodeError := json.NewDecoder(r.Body).Decode(update); decodeError != nil {
		log.Printf("No se ha podido decodificar el contenido de la petición: %s", decodeError)
		return
	}

	var message string

	if update.ChannelPost != nil {
		log.Printf("[Channel - %s] %s", update.ChannelPost.SenderChat.Title, update.ChannelPost.Text)
		message = update.ChannelPost.Text
	} else if update.Message != nil {
		log.Printf("[User - %s] %s", update.Message.From.UserName, update.Message.Text)
		message = update.Message.Text
	}

	// TODO: Mandar mensaje a la API

	_, printError := fmt.Fprint(w, "Se ha enviado el mensaje a la API. Contenido: "+message)
	if printError != nil {
		return
	}
}

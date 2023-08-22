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

	listenPort := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenPort = ":" + val
	}

	http.HandleFunc("/api/setup", func(writer http.ResponseWriter, request *http.Request) {
		setUpWebhook(writer, request, bot)
	})
	http.HandleFunc("/api/handleUpdate", handleTelegramUpdate)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}

func setUpWebhook(writer http.ResponseWriter, response *http.Request, bot *telegram.BotAPI) {

	webHook, _ := telegram.NewWebhook("https://" + response.Host + "/api/handleUpdate")

	_, webHookError := bot.Request(webHook)

	if webHookError != nil {
		log.Fatal(webHookError)
	}

	_, printError := fmt.Fprint(writer, "Se ha configurado el webhook correctamente")
	if printError != nil {
		return
	}
}

func handleTelegramUpdate(writer http.ResponseWriter, request *http.Request) {

	update := &telegram.Update{}
	if decodeError := json.NewDecoder(request.Body).Decode(update); decodeError != nil {
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

	_, printError := fmt.Fprint(writer, "Se ha enviado el mensaje a la API. Contenido: "+message)
	if printError != nil {
		return
	}
}

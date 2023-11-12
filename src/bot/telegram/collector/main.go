package main

import (
	"log"
	"os"
	"vidya-sale/bot/telegram/collector/handlers"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	listenPort := ":" + os.Getenv("PORT")

	app.Post("/collect", handlers.CollectUpdate)

	setUpWebHook()
	app.Listen(listenPort)
}

func setUpWebHook() {
	bot, errorBot := telegram.NewBotAPI(os.Getenv("BOT_TOKEN"))

	if errorBot != nil {
		log.Fatal(errorBot)
	}

	webHook, urlError := telegram.NewWebhook(os.Getenv("HOST_URL") + "/collect")

	if urlError != nil {
		log.Fatal(urlError)
	}

	_, setWebHookError := bot.Request(webHook)

	if setWebHookError != nil {
		log.Fatal(setWebHookError)
	}
}

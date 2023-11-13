package main

import (
	"log"
	"os"
	"vidya-sale/bot/telegram/messenger/handlers"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	bot := setUpBot()

	log.Print("[Bot] Bot is up and running")

	app.Post("/fetch-offers", func(context *fiber.Ctx) error {
		return handlers.GetGameOffers(context, bot)
	})

	setUpWebHook(bot)
	app.Listen(":" + os.Getenv("PORT"))
}

func setUpBot() *telegram.BotAPI {
	bot, errorBot := telegram.NewBotAPI(os.Getenv("BOT_TOKEN"))

	if errorBot != nil {
		log.Fatal(errorBot)
	}

	return bot
}

func setUpWebHook(bot *telegram.BotAPI) {
	webHook, urlError := telegram.NewWebhook(os.Getenv("HOST_URL") + "/fetch-offers")

	if urlError != nil {
		log.Fatal(urlError)
	}

	log.Printf("[Bot] Webhook created successfully for address %s", webHook.URL.String())

	_, setWebHookError := bot.Request(webHook)

	if setWebHookError != nil {
		log.Fatal(setWebHookError)
	}

	log.Print("[Bot] Webhook set successfully")
}

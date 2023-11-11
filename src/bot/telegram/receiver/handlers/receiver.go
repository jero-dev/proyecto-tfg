package handlers

import (
	"bytes"
	"log"
	"net/http"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
)

var (
	MessageNotAbleToParseRequestBody = "Could not parse the request body"
	MessageCouldNotSendToAPI         = "Could not send the message to the API"
	MessageOK                        = "Process finished correctly"
)

func HandleUpdates(context *fiber.Ctx) error {
	update := &telegram.Update{}

	if parseError := context.BodyParser(update); parseError != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": MessageNotAbleToParseRequestBody,
		})
	}

	if update.ChannelPost != nil {
		log.Printf("[Channel - %s] %s", update.ChannelPost.SenderChat.Title, update.ChannelPost.Text)
		message := update.ChannelPost.Text
		apiURL := os.Getenv("API_URL")

		_, responseError := http.Post(apiURL, "application/json", bytes.NewBuffer([]byte(message)))

		if responseError != nil {
			return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": MessageCouldNotSendToAPI,
			})
		}
	}

	return context.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": MessageOK,
	})
}

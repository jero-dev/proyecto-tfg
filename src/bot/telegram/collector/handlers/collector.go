package handlers

import (
	"bytes"
	"encoding/json"
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

func CollectUpdate(context *fiber.Ctx) error {
	update := &telegram.Update{}

	if parseError := context.BodyParser(update); parseError != nil {
		log.Printf("Couldn't parse the request. Error: %s", parseError.Error())
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": MessageNotAbleToParseRequestBody,
		})
	}

	if update.ChannelPost != nil {
		log.Printf("[Channel - %s] %s", update.ChannelPost.SenderChat.Title, update.ChannelPost.Text)
		message := struct {
			Message string `json:"message"`
		}{
			Message: update.ChannelPost.Text,
		}
		body, _ := json.Marshal(message)
		apiURL := os.Getenv("API_URL") + "/offers"

		response, responseError := http.Post(apiURL, "application/json", bytes.NewBuffer(body))

		if responseError != nil {
			log.Printf("Couldn't send the message to the API. Error: %s", responseError.Error())
			return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": MessageCouldNotSendToAPI,
			})
		}

		if response.StatusCode != http.StatusOK {
			log.Printf("The API returned a %s status code.", response.Status)
			return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "The API returned a " + response.Status + " status code.",
			})
		}
	}

	log.Println("Message stored successfully.")
	return context.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": MessageOK,
	})
}

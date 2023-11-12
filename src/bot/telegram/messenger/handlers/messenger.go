package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
)

var (
	MessageNotAbleToParseRequestBody    = "Could not parse the request body"
	MessageCouldNotSendToAPI            = "Could not retrieve the information from the API"
	MessageCouldNotSendAnswerToTelegram = "Could not send the answer to Telegram"
	MessageOK                           = "Process finished correctly"
)

type GameOffersResponse struct {
	Platforms []PlatformResponse `json:"data"`
}

type PlatformResponse struct {
	Platform string          `json:"platform"`
	Offers   []OfferResponse `json:"offers"`
}

type OfferResponse struct {
	Link  string `json:"link"`
	Price string `json:"price"`
}

func GetGameOffers(context *fiber.Ctx, bot *telegram.BotAPI) error {
	update := &telegram.Update{}

	if parseError := context.BodyParser(update); parseError != nil {
		log.Printf("[Error parsing the request body] %s", parseError.Error())
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": MessageNotAbleToParseRequestBody,
		})
	}

	if update.Message != nil && update.Message.Text != "" {
		log.Printf("[Message received] %s", update.Message.Text)
		message := update.Message.Text
		apiURL := os.Getenv("API_URL") + "/offers/" + url.QueryEscape(message)

		gameOffersResponse, responseError := http.Get(apiURL)

		if responseError != nil {
			log.Printf("[Error sending the request to the API] %s", responseError.Error())
			return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": MessageCouldNotSendToAPI,
			})
		}

		if gameOffersResponse.StatusCode != http.StatusOK &&
			gameOffersResponse.StatusCode != http.StatusNotFound {
			log.Printf("[Not appealing API response] %s", gameOffersResponse.Status)
			return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "The API returned an unexpected status code: " + gameOffersResponse.Status,
			})
		}

		defer gameOffersResponse.Body.Close()
		answerMessage := composeMessage(gameOffersResponse.StatusCode,
			gameOffersResponse.Body, update.Message.Chat.ID)

		if _, errorBot := bot.Send(answerMessage); errorBot != nil {
			log.Printf("[Error sending the answer to Telegram] %s", errorBot.Error())
			return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": MessageCouldNotSendToAPI,
			})
		}
	}

	log.Print("Process finished correctly.")
	return context.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": MessageOK,
	})
}

func composeMessage(statusCode int, body io.ReadCloser, chatID int64) telegram.MessageConfig {
	var answerMessage telegram.MessageConfig

	if statusCode == http.StatusNotFound {
		answerMessage = telegram.NewMessage(chatID,
			"No he encontrado ninguna oferta para ese juego.")
	}

	if statusCode == http.StatusOK {
		gameOffers := &GameOffersResponse{}
		json.NewDecoder(body).Decode(gameOffers)

		answer := "Estas son las ofertas que he encontrado: \n\n"
		for _, platform := range gameOffers.Platforms {
			answer += platform.Platform + ":\n"
			for _, offer := range platform.Offers {
				answer += offer.Link + " - " + offer.Price + "€\n"
			}
		}

		answerMessage = telegram.NewMessage(chatID, answer)
	}

	return answerMessage
}

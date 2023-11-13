package main

import (
	"net/url"
	services "vidya-sale/api/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	api := app.Group("/api", logger.New())
	listenPort := ":8080"

	managerService, serviceError :=
		services.NewOfferManagerService(services.WithPostgresProductRepository())

	if serviceError != nil {
		panic(serviceError)
	}

	processorService := services.NewMessageProcessorService()

	api.Get("/offers/:gameName", func(context *fiber.Ctx) error {
		return getGameOffers(context, managerService)
	})
	api.Post("/offers", func(context *fiber.Ctx) error {
		return storeOffer(context, managerService, processorService)
	})

	app.Listen(listenPort)
}

func getGameOffers(context *fiber.Ctx, offerManager *services.OfferManagerService) error {
	gameName, _ := url.QueryUnescape(context.Params("gameName"))

	gameOffers, getGameOffersError := offerManager.GetGameOffers(gameName)

	if getGameOffersError != nil {
		return context.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Could not find offers for the related game.",
		})
	}

	offersResponse := []PlatformResponse{}
	for platform, offers := range gameOffers {
		platformOffers := []OfferResponse{}
		for _, offer := range offers {
			platformOffers = append(platformOffers, OfferResponse{
				Link:  offer.GetLink(),
				Price: offer.GetPrice(),
			})
		}
		offersResponse = append(offersResponse, PlatformResponse{
			Platform: platform,
			Offers:   platformOffers,
		})
	}

	return context.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Game offers successfully retrieved.",
		"data":    offersResponse,
	})
}

func storeOffer(context *fiber.Ctx, offerManager *services.OfferManagerService,
	processorService *services.MessageProcessorService) error {

	messageRequest := &MessageRequest{}

	if parseError := context.BodyParser(messageRequest); parseError != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Could not parse request body.",
		})
	}

	gameName, platform, link, price := processorService.ParseMessage(messageRequest.Message)

	storingError := offerManager.StoreOffer(gameName, platform, link, price)

	if storingError != nil {
		return context.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Could not find the elements correctly.",
			"data": fiber.Map{
				"incomingMessage": messageRequest.Message,
				"gameName":        gameName,
				"platform":        platform,
				"link":            link,
				"price":           price,
			},
		})
	}

	return context.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Offer successfully stored.",
	})
}

type MessageRequest struct {
	Message string `json:"message"`
}

type PlatformResponse struct {
	Platform string          `json:"platform"`
	Offers   []OfferResponse `json:"offers"`
}

type OfferResponse struct {
	Link  string  `json:"link"`
	Price float64 `json:"price"`
}

package main

import (
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	aggregates "vidya-sale/api/aggregate"
	"vidya-sale/api/domain/product/memory"
	services "vidya-sale/api/service"
	valueobjects "vidya-sale/api/valueobject"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetOffersRoute(t *testing.T) {
	app := fiber.New()
	existingProduct, _ := aggregates.NewProduct("Red Dead Redemption", "Switch")
	existingOffer, _ := valueobjects.NewOffer("Test Link", 44.99)
	existingProduct.AddOffer(existingOffer)
	productRepository := memory.New()
	productRepository.Add(existingProduct)
	offerManager, _ := services.NewOfferManagerService(services.WithProductRepository(productRepository))

	app.Get("/api/offers/:gameName", func(context *fiber.Ctx) error {
		return getGameOffers(context, offerManager)
	})

	testCases := []struct {
		name         string
		route        string
		expectedCode int
	}{
		{
			name:         "Get HTTP status 200 when offers exist",
			route:        "/api/offers/" + url.QueryEscape(existingProduct.GetName()),
			expectedCode: fiber.StatusOK,
		},
		{
			name:         "Get HTTP status 404 when offers do not exist",
			route:        "/api/offers/" + url.QueryEscape("Non existing game"),
			expectedCode: fiber.StatusNotFound,
		},
		{
			name:         "Get HTTP status 404 when route does not exist",
			route:        "/not-found",
			expectedCode: fiber.StatusNotFound,
		},
	}

	for _, testCase := range testCases {
		request := httptest.NewRequest("GET", testCase.route, nil)

		response, _ := app.Test(request, 1)

		assert.Equalf(t, testCase.expectedCode, response.StatusCode, testCase.name)
	}
}

func TestStoreOfferRoute(t *testing.T) {
	app := fiber.New()
	productRepository := memory.New()
	offerManager, _ := services.NewOfferManagerService(services.WithProductRepository(productRepository))
	processorService := services.NewMessageProcessorService()

	properMessageBody := struct {
		Message string `json:"message,omitempty"`
	}{
		Message: `📆⬇️ Star Ocean The Second Story R #Switch
		(https://telegra.ph/file/2c7ec774526506600df7a.png)BAJONAZO FLASH a solo 47,51€, precio mínimo alcanzado incluso sin cupón (PVP 60€)
	   ✂️ Con cupón de primera compra se te queda a solo 33,25€
	   
	   🌸 https://ojueg.es/vgfPa
	   
	   
	   📝 Edición española con envío de lanzamiento de parte de Gamers4Life
	   
	   
	   ⚪️ Si no te aclaras con el tema del cupón, averigua como sacarle partido en nuestro hilo de dudas Miravia del @GrupoNintendoOJ 
	   🍨 https://t.me/GrupoNintendoOJ/481642
		`,
	}
	goodBodyMessage, _ := json.Marshal(properMessageBody)

	app.Post("/api/offers", func(context *fiber.Ctx) error {
		return storeOffer(context, offerManager, processorService)
	})

	testCases := []struct {
		name         string
		route        string
		body         string
		expectedCode int
	}{
		{
			name:         "Get HTTP status 200 when offer is stored",
			route:        "/api/offers",
			body:         string(goodBodyMessage),
			expectedCode: fiber.StatusOK,
		},
		{
			name:         "Get HTTP status 400 when body is not parsed",
			route:        "/api/offers",
			body:         `{"message": "Test Message"`,
			expectedCode: fiber.StatusBadRequest,
		},
		{
			name:         "Get HTTP status 422 when message is unprocessable",
			route:        "/api/offers",
			body:         `{"message": "Test Message"}`,
			expectedCode: fiber.StatusUnprocessableEntity,
		},
	}

	for _, testCase := range testCases {
		request := httptest.NewRequest("POST", testCase.route, strings.NewReader(testCase.body))
		request.Header.Set("Content-Type", "application/json")

		response, _ := app.Test(request, 1)

		assert.Equalf(t, testCase.expectedCode, response.StatusCode, testCase.name)
	}
}

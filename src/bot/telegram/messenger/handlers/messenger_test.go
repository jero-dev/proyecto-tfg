package handlers

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func Test_GetGameOffers(t *testing.T) {
	app := fiber.New()

	app.Get("/fetch-offers", func(context *fiber.Ctx) error {
		return GetGameOffers(context, nil)
	})

	testCases := []struct {
		name         string
		route        string
		body         string
		expectedCode int
	}{
		{
			name:         "Get HTTP status 200 when process finished correctly",
			route:        "/fetch-offers",
			body:         `{"update_id": 2, "channel_post": {"text": "Text Message", "sender_chat": {"title": "Test Channel"}}}`,
			expectedCode: fiber.StatusOK,
		},
		{
			name:         "Get HTTP status 500 when message could not be sent to the API",
			route:        "/fetch-offers",
			body:         `{"update_id": 1, "message": {"text": "Test Game"}}`,
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:         "Get HTTP status 400 when message is unparsable",
			route:        "/fetch-offers",
			body:         `{"invalid": }`,
			expectedCode: fiber.StatusBadRequest,
		},
		{
			name:         "Get HTTP status 404 when route does not exist",
			route:        "/not-found",
			body:         `{}`,
			expectedCode: fiber.StatusNotFound,
		},
	}

	for _, testCase := range testCases {
		request := httptest.NewRequest("GET", testCase.route, bytes.NewBufferString(testCase.body))
		request.Header.Set("Content-Type", "application/json")

		response, _ := app.Test(request, -1)

		assert.Equalf(t, testCase.expectedCode, response.StatusCode, testCase.name)
	}
}

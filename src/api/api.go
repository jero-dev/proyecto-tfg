package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	api := app.Group("/api")

	api.Get("/hello", func(context *fiber.Ctx) error {
		return context.SendString("Hello, World!")
	})

	apiError := app.Listen(":3000")

	if apiError != nil {
		panic(apiError)
	}
}

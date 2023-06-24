package main

import (
	"os"
	"projector/queue"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}

	go queue.ConsumeEvents()

	app.Listen(":" + port)
}

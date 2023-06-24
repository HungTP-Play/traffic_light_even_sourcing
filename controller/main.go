package main

import (
	"controller/queue"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "2222"
	}

	go queue.ConsumeEvents()
	app.Listen(":" + port)
}

package main

import (
	"os"
	"projector/metrics"
	"projector/queue"
	"projector/repo"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	defer repo.CloseDB()

	app.Use(metrics.PrometheusMiddleware)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/metrics", metrics.Metrics)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}

	go queue.ConsumeEvents()

	app.Listen(":" + port)
}

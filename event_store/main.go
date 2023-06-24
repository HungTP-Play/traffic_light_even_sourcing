package main

import (
	"encoding/json"
	"event_store/model"
	"event_store/queue"
	"event_store/repo"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func dispatch(event model.EventEmitDto) error {
	// Get event from event_in
	return queue.PublishEvent(event)
}

func main() {
	app := fiber.New()
	defer repo.CloseDB()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	app.Post("/emit", func(c *fiber.Ctx) error {
		body := c.Body()

		// Marshal body to map[string]interface{}
		event := model.EventEmitDto{}
		err := json.Unmarshal(body, &event)
		if err != nil {
			// Return with 400
			return c.Status(400).JSON(fiber.Map{
				"status": "failed",
				"error":  err.Error(),
			})
		}

		err = repo.StoreEvent(event)
		if err != nil {
			// Return with 500
			return c.Status(500).JSON(fiber.Map{
				"status": "failed",
				"error":  err.Error(),
			})
		}

		// Dispatch event to queue

		err = dispatch(event)
		if err != nil {
			// Return with 500
			return c.Status(500).JSON(fiber.Map{
				"status": "failed",
				"error":  err.Error(),
			})
		}
		return c.JSON(map[string]interface{}{
			"status": "success",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "1111"
	}

	app.Listen(":" + port)
}

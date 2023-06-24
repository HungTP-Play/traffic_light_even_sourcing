package main

import (
	"encoding/json"
	"event_store/model"
	"event_store/queue"
	"event_store/repo"
	"os"

	"github.com/gofiber/fiber/v2"
)

func dispatch(eventIn chan model.EventEmitDto, eventOut chan error) {
	// Get event from event_in
	event := <-eventIn
	err := queue.PublishEvent(event)
	eventOut <- err
}

func main() {
	app := fiber.New()

	app.Post("/emit", func(c *fiber.Ctx) error {
		eventIn := make(chan model.EventEmitDto)
		eventOut := make(chan error)
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

		go dispatch(eventIn, eventOut)
		eventResponse := <-eventOut
		if eventResponse != nil {
			// Return with 500
			return c.Status(500).JSON(fiber.Map{
				"status": "failed",
				"error":  eventResponse.Error(),
			})
		}
		return c.JSON(eventResponse)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "1111"
	}

	app.Listen(":" + port)
}

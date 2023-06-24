package queue

import (
	"controller/model"
	"controller/repo"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func HandleMessage(event string) {
	fmt.Printf("[[Controller]] Received event: %s\n", event)
	var eventEmitDto model.EventEmitDto
	err := json.Unmarshal([]byte(event), &eventEmitDto)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal event: %s", err)
	}

	if eventEmitDto.EventName == "registration_event" {
		newLight := model.TrafficLight{
			LightID:      eventEmitDto.EventData.(map[string]interface{})["light_id"].(string),
			Location:     eventEmitDto.EventData.(map[string]interface{})["location"].(string),
			RegisteredAt: time.Now().Unix(),
		}
		err = repo.StoreTrafficLight(newLight)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to store traffic light: %s", err)
		}
	}
}

func ConsumeEvents() error {
	// Wait for 30 secs
	time.Sleep(30 * time.Second)

	connString := os.Getenv("RABBITMQ_CONNECTION_STRING")
	if connString == "" {
		fmt.Fprintln(os.Stderr, "RABBITMQ_CONNECTION_STRING is not set")
	}

	conn, err := amqp.Dial(connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	controllerQueueName := os.Getenv("CONTROLLER_QUEUE")
	q, err := ch.QueueDeclare(
		controllerQueueName, // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}
	fmt.Printf("[[Controller]] Start consuming events from queue: %s 🎉🎉\n", controllerQueueName)

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			event := string(msg.Body)
			HandleMessage(event)
		}
	}()

	<-forever

	return nil
}

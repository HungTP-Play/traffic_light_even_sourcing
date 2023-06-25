package queue

import (
	"encoding/json"
	"fmt"
	"os"
	"projector/metrics"
	"projector/model"
	"projector/repo"
	"projector/util"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection

func init() {
	time.Sleep(10 * time.Second)
	conn = Connect()
}

func colorToInt(color string) int {
	if color == "RED" {
		return 1
	} else if color == "GREEN" {
		return 2
	} else if color == "YELLOW" {
		return 3
	} else {
		return 0
	}
}

func HandleMessage(event string) {
	var eventEmitDto model.EventEmitDto
	err := json.Unmarshal([]byte(event), &eventEmitDto)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal event: %s", err)
	}

	if eventEmitDto.EventName == "registration_event" {
		location := eventEmitDto.EventData.(map[string]interface{})["location"].(string)
		// The location is in the format of "lat::lng"
		locationLat := location[:strings.Index(location, "::")]
		locationLng := location[strings.Index(location, "::")+2:]
		lightID := eventEmitDto.EventData.(map[string]interface{})["light_id"].(string)
		color := "RED"

		// And default color is RED=1, GREEN=2, YELLOW=3
		metrics.SetTrafficLightState(lightID, locationLat, locationLng, util.StringColorToInt(color))

		repo.UpsertTrafficLight(lightID, util.StringToFloat64(locationLat), util.StringToFloat64(locationLng), util.StringColorToInt(color))
	}

	if eventEmitDto.EventName == "state_change_event" {
		location := eventEmitDto.EventData.(map[string]interface{})["location"].(string)
		locationLat := location[:strings.Index(location, "::")]
		locationLng := location[strings.Index(location, "::")+2:]
		lightID := eventEmitDto.EventData.(map[string]interface{})["light_id"].(string)
		color := eventEmitDto.EventData.(map[string]interface{})["to_state"].(string)
		// And default color is RED=1, GREEN=2, YELLOW=3
		metrics.SetTrafficLightState(lightID, locationLat, locationLng, util.StringColorToInt(color))

		repo.UpsertTrafficLight(lightID, util.StringToFloat64(locationLat), util.StringToFloat64(locationLng), util.StringColorToInt(color))
	}
}

func Connect() *amqp.Connection {
	connString := os.Getenv("RABBITMQ_CONNECTION_STRING")
	if connString == "" {
		fmt.Fprintln(os.Stderr, "RABBITMQ_CONNECTION_STRING is not set")
	}

	conn, err := amqp.Dial(connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to RabbitMQ: %s", err)
	}

	return conn
}

func Close() {
	conn.Close()
}

func GetConnection() *amqp.Connection {
	if conn == nil {
		conn = Connect()
	}

	return conn
}

func ConsumeEvents() error {
	time.Sleep(15 * time.Second)
	conn := GetConnection()
	if conn.IsClosed() {
		conn = Connect()
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	projectionQueueName := os.Getenv("PROJECTION_QUEUE")
	q, err := ch.QueueDeclare(
		projectionQueueName, // name
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
	fmt.Printf("[[Projector]] Start consuming events from queue: %s 🎉🎉\n", projectionQueueName)

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

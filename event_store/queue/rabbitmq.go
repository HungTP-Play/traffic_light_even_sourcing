package queue

import (
	"encoding/json"
	"event_store/model"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func getRabbitMQConnection() string {
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitPort := os.Getenv("RABBITMQ_PORT")
	rabbitUser := "guest"
	rabbitPassword := "guest"
	return "amqp://" + rabbitUser + ":" + rabbitPassword + "@" + rabbitHost + ":" + rabbitPort + "/"
}

var conn *amqp.Connection

func Connect() *amqp.Connection {
	connectionString := getRabbitMQConnection()
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		panic("[[Event Store]] failed to connect to RabbitMQ")
	}
	return conn
}

func Close() {
	conn.Close()
}

func GetRabbitMQConnection() *amqp.Connection {
	if conn == nil {
		conn = Connect()
	}

	return conn
}

func PublishProjectionEvents(event model.EventEmitDto) error {
	conn := GetRabbitMQConnection()
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

	eventMsg, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventMsg,
		},
	)
	if err != nil {
		return err
	}
	fmt.Printf("[[Event Store]] Published event: %v to %v 🎉\n", event, projectionQueueName)
	return nil
}

func PublishControllerEvents(event model.EventEmitDto) error {
	conn := GetRabbitMQConnection()
	if conn.IsClosed() {
		conn = Connect()
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

	eventMsg, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventMsg,
		},
	)
	if err != nil {
		return err
	}

	fmt.Printf("[[Event Store]] Published event: %v to %v 🎉\n", event, controllerQueueName)
	return nil
}

func PublishEvent(event model.EventEmitDto) error {
	switch event.EventName {
	case "registration_event":
		err := PublishControllerEvents(event)
		if err != nil {
			return err
		}

		err = PublishProjectionEvents(event)
		if err != nil {
			return err
		}
		fmt.Println("[[Event Store]] published event: registration to metadata and projections 🎉🎉")
		return nil
	case "state_change_event":
		return PublishProjectionEvents(event)
	case "light_state_override_event":
		return PublishProjectionEvents(event)
	case "light_state_override_done_event":
		return PublishProjectionEvents(event)
	default:
		return nil
	}
}

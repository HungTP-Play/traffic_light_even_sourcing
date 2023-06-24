package queue

import (
	"encoding/json"
	"event_store/model"
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
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"projections", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
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
	return nil
}

func PublishMetadataEvents(event model.EventEmitDto) error {
	conn := GetRabbitMQConnection()
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"metadata", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
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
	return nil
}

func PublishEvent(event model.EventEmitDto) error {
	switch event.EventName {
	case "registration":
		err := PublishMetadataEvents(event)
		if err != nil {
			return err
		}

		err = PublishProjectionEvents(event)
		if err != nil {
			return err
		}
		return nil
	case "state_change":
		return PublishProjectionEvents(event)
	case "light_state_override":
		return PublishProjectionEvents(event)
	case "light_state_override_done":
		return PublishProjectionEvents(event)
	default:
		return nil
	}
}

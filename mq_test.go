package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"testing"
)

func TestPushMsg(t *testing.T) {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
		}
	}(ch)

	for true {
		err = ch.PublishWithContext(context.Background(), "", "a", // routing key
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("1234"),
			})
		failOnError(err, "Failed to publish a message")
	}

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

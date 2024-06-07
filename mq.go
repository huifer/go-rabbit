package main

// RMQ PACKAGE - "rmq"
import (
	"errors"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rmqCredentials string = "amqp://guest:guest@localhost:5672"
	rmqContentType string = "text/plain"
)

var conn *amqp.Connection
var ch *amqp.Channel

func hasError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func ConnectToRMQ() (err error) {
	conn, err = amqp.Dial(rmqCredentials)
	if err != nil {
		return errors.New("Error de conexion: " + err.Error())
	}

	ch, err = conn.Channel()
	if err != nil {
		return errors.New("Error al abrir canal: " + err.Error())
	}

	q, err := ch.QueueDeclare("a", true, false, false, false, nil)
	if err != nil {
		return errors.New("Error al abrir canal: " + err.Error())
	}
	_ = q

	ch.Qos(1, 0, false)
	if err != nil {
		return errors.New("Error al abrir canal: " + err.Error())
	}

	observeConnection()

	return nil
}

func observeConnection() {
	go func() {
		log.Printf("Conexion perdida: %s\n", <-conn.NotifyClose(make(chan *amqp.Error)))
		log.Printf("Intentando reconectar con RMQ\n")

		closeActiveConnections()

		for err := ConnectToRMQ(); err != nil; err = ConnectToRMQ() {
			log.Println(err)
			time.Sleep(5 * time.Second)
		}
	}()
}

// Can be also implemented in graceful shutdowns
func closeActiveConnections() {
	if !ch.IsClosed() {
		if err := ch.Close(); err != nil {
			log.Println(err.Error())
		}
	}

	if conn != nil && !conn.IsClosed() {
		if err := conn.Close(); err != nil {
			log.Println(err.Error())
		}
	}
}

// SendMessage - message without response
func SendMessage(body string) {
	err := ch.Publish(
		"",    // exchange
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: rmqContentType,
			Body:        []byte(body),
		})

	if err != nil {
		log.Printf("%s\n %s\n", "Error al publicar mensaje", err)
		log.Println(body)
	}
}

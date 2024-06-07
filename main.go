package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	ConnectToRMQ()
	preHandlerMessage, _ := chann.Consume("a", // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	go startHttp()
	for {
		select {
		case msg := <-preHandlerMessage:
			fmt.Println("Received a message: ", string(msg.Body))
			msg.Ack(false)
		}
	}
}

func startHttp() {
	http.ListenAndServe(":19000", nil)

}

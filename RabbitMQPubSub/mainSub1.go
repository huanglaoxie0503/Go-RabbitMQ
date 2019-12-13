package main

import "Go-RabbitMQ/RabbitMQHelper"

func main() {
	rabbitMq := RabbitMQHelper.NewRabbitMQPubSub("" + "newProduct")
	rabbitMq.ReceiveSub()
}

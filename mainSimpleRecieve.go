package main

import "Go-RabbitMQ/RabbitMQHelper"

func main() {
	rabbitMq := RabbitMQHelper.NewRabbitMQSimple("" + "stockSimple")
	rabbitMq.ConsumeSimple()
}

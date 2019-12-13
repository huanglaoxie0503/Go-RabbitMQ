package main

import "Go-RabbitMQ/RabbitMQHelper"

func main() {
	stockOne := RabbitMQHelper.NewRabbitMQTopic("exStockTopic", "#")
	stockOne.ReceiveTopic()
}

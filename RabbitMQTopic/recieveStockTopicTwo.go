package main

import "Go-RabbitMQ/RabbitMQHelper"

func main() {
	stockOne := RabbitMQHelper.NewRabbitMQTopic("exStockTopic", "stock.*.two")
	stockOne.ReceiveTopic()
}

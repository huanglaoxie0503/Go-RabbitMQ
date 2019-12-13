package main

import "Go-RabbitMQ/RabbitMQHelper"

func main() {
	stockOne := RabbitMQHelper.NewRabbitMQRouting("exStock", "stock_one")
	stockOne.ReceiveRouting()
}

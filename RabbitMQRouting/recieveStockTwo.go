package main

import "Go-RabbitMQ/RabbitMQHelper"

func main() {
	stockTwo := RabbitMQHelper.NewRabbitMQRouting("exStock", "stock_one")
	stockTwo.ReceiveRouting()
}

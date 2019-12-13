package main

import (
	"Go-RabbitMQ/RabbitMQHelper"
	"fmt"
)

func main() {
	rabbitMq := RabbitMQHelper.NewRabbitMQSimple("" + "stockSimple")
	rabbitMq.PublishSimple("我要买贵州茅台10000股")
	fmt.Println("发送成功！")
}

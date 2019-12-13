package main

import (
	"Go-RabbitMQ/RabbitMQHelper"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitMq := RabbitMQHelper.NewRabbitMQPubSub("" + "newProduct")
	for i := 0; i < 100; i ++ {
		rabbitMq.PublishPub("订阅模式生产第" + strconv.Itoa(i) + "条" + "数据")
		fmt.Println("订阅模式生产第" + strconv.Itoa(i) + "条" + "数据")
		time.Sleep(1*time.Second)
	}
}

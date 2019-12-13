package main

import (
	"Go-RabbitMQ/RabbitMQHelper"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitMq := RabbitMQHelper.NewRabbitMQSimple("" + "stockSimple")
	for i := 0; i <= 100 ; i ++  {
		rabbitMq.PublishSimple("Hello 贵州茅台" + strconv.Itoa(i))
		time.Sleep(1*time.Second)
		fmt.Println(i)
	}
}

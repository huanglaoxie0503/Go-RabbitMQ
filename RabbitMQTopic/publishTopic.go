package main

import (
	"Go-RabbitMQ/RabbitMQHelper"
	"fmt"
	"strconv"
	"time"
)

func main() {
	stockOne := RabbitMQHelper.NewRabbitMQTopic("exStockTopic", "stock.topic.one")
	stockTwo := RabbitMQHelper.NewRabbitMQTopic("exStockTopic", "stock.topic.two")
	for i := 0; i <= 10 ; i ++  {
		stockOne.PublishTopic("我要买贵州茅台！！！！" + strconv.Itoa(i))
		stockTwo.PublishTopic("我要买贵州茅台！！！！" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}

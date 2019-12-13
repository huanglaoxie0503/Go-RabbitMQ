package main

import (
	"Go-RabbitMQ/RabbitMQHelper"
	"fmt"
	"strconv"
)

func main() {
	stockOne := RabbitMQHelper.NewRabbitMQRouting("exStock", "stock_one")
	stockTwo := RabbitMQHelper.NewRabbitMQRouting("exStock", "stock_two")
	for i := 0;i<=10 ;i++  {
		stockOne.PublishRouting("我要买贵州茅台！！！！" + strconv.Itoa(i))
		stockTwo.PublishRouting("我要买贵州茅台！！！！" + strconv.Itoa(i))
		fmt.Println(i)
	}
}

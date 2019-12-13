package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

// 声明变量
var conn *amqp.Connection // rabbitMq 连接
var channel *amqp.Channel
var count = 0

const (
	queueName = "stock"
	exchange = ""
	mQUrl = "amqp://guest:guest@127.0.0.1:5673"
)

func failOnErr(err error, msg string)  {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
		panic(fmt.Sprintf("%s:%s", msg, err))
	}
}

// rabbitMq 连接函数
func Connect()  {
	var err error
	// 连接 rabbitMq
	conn, err := amqp.Dial(mQUrl)
	failOnErr(err, "failed to connect")
	// 获取 channel
	channel, err = conn.Channel()
	failOnErr(err, "failed to open a channel")
}

// 关闭rabbitMq连接
func Close()  {
	// 1. 关闭channel
	_ = channel.Close()
	// 2. 关闭连接
	_ = conn.Close()
}

func Push()  {
	// 1. 判断是否存在channel
	if channel == nil {
		Connect()
	}
	// 2. 消息
	message := "哈喽，贵州茅台!"
	// 3. 声明队列
	q, err := channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
		)
	// 4. 判断错误
	if err != nil {
		fmt.Println(err)
	}

	// 5. 生产消息
	_ = channel.Publish(exchange, q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body: []byte(message),
	})
}

// 消费端
func Receive()  {
	// 1. 判断 channel 是否存在
	if channel == nil {
		Connect()
	}

	// 2. 声明队列
	q, err := channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
		)
	//3.消费代码
	msg, err := channel.Consume(q.Name, "",
		true,
		false,
		false,
		false,
		nil)
	failOnErr(err, "获取消费异常信息")

	msgForever := make(chan bool)

	// 消费逻辑
	go func() {
		for d := range msg {
			s := string(d.Body)
			count ++
			fmt.Printf("接收信息是%s-- %d\n", s, count)
		}
	}()
	fmt.Printf("退出请按 CTRL+C\n")

	<-msgForever
}


func main() {
	go func() {
		for {
			Push()
			time.Sleep(5 * time.Second)
		}
	}()
	Receive()
	fmt.Println("生产消费完成")
	Close()
}

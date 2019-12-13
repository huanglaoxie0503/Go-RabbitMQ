package RabbitMQHelper

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// 连接信息
const URL  = "amqp://oscar:oscar@127.0.0.1:5672/news"

type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机名称
	Exchange string
	//bind key 名称
	key string
	// 连接信息
	MqUrl string
}

// 创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{
		conn:      nil,
		channel:   nil,
		QueueName: queueName,
		Exchange:  exchange,
		key:       key,
		MqUrl:     URL,
	}
}

// 断开channel 和 connection
func (r *RabbitMQ) Destroy() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}

// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// 创建简单模式下打RabbitMQ 实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	// 创建rabbitMQ 实例
	rabbitMq := NewRabbitMQ(queueName, "", "")
	var err error

	// 获取connection
	rabbitMq.conn, err = amqp.Dial(rabbitMq.MqUrl)
	rabbitMq.failOnErr(err, "failed to connect rabbitMq")

	// 获取channel
	rabbitMq.channel, err = rabbitMq.conn.Channel()
	rabbitMq.failOnErr(err, "failed to open a channel")

	return rabbitMq
}

// 直接模式队列生产
func (r *RabbitMQ) PublishSimple(message string) {
	// 1. 申请队列，如果队列不存在则自动创建，存在则跳过
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 是否自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞处理
		false,
		// 额外的属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	// 调用channel 发送消息到队列中
	_ = r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// 如果为true， 根据自身exchange 类型和 routekey 规则无法找到符合条件打队列会把消息返还给发送者
		false,
		// 如果为true, 当exchange 发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// Simple 模式下打消费者
func (r *RabbitMQ) ConsumeSimple() {
	// 1. 申请队列， 如果对列不存在自动创建，存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 是否自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞处理
		false,
		// 额外属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	// 接受消息
	msgAll, err := r.channel.Consume(
		// queue
		q.Name,
		// consumer
		"",
		// auto-ack
		true,
		// exclusive
		false,
		// 设置为true， 表示 不能将同一个connection 中生产者发送打消息传递给这个Connection 中打消费者
		false,
		// 队列是否阻塞
		false,
		// 额外属性 args
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	// 启动协程处理消息
	go func() {
		for d := range msgAll{
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<- forever
}

// 订阅模式创建RabbitMQ 实例

func NewRabbitMQPubSub(exchange string) *RabbitMQ {
	// 创建 RabbitMQHelper 实例
	rabbitMq := NewRabbitMQ("", exchange, "")
	var err error

	// 获取connection
	rabbitMq.conn, err = amqp.Dial(rabbitMq.MqUrl)
	rabbitMq.failOnErr(err, "failed to connect rabbitMq")

	// 获取channel
	rabbitMq.channel, err = rabbitMq.conn.Channel()
	rabbitMq.failOnErr(err, "failed to open a channel")

	return rabbitMq
}

// 订阅生产模式
func (r *RabbitMQ)PublishPub(message string) {
	// 1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		// true 表示 exchange 不可以被client 用来推送消息，仅用来进行exchange 和 exchange 之间打绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	// 2. 发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

}

// 订阅模式消费代码
func (r *RabbitMQ) ReceiveSub() {
	// 1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 交换机类型
		"fanout",
		true,
		false,
		// true 表示这个exchange 不可以被client 用来推送消息，仅用来进行 exchange 和 exchange 之间打绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	// 2. 尝试创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		// 随机生产队列名称
		"",
		false,
		false,
		true,
		false,
		nil,
		)
	r.failOnErr(err, "Failed to declare a queue")

	// 绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		"",
		r.Exchange,
		false,
		nil,
		)

	// 消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
		)
	forever := make(chan bool)

	go func() {
		for d:= range messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	fmt.Println("退出请按 CTRL+C\n")
	<- forever
}

// 路由模式 创建RabbitMQ 实例
func NewRabbitMQRouting(exchangeName string, routingKey string)*RabbitMQ {
	// 创建RabbitMQ 实例
	rabbitMq := NewRabbitMQ("", exchangeName, routingKey)
	var err error

	// 获取connection
	rabbitMq.conn,err = amqp.Dial(rabbitMq.MqUrl)
	rabbitMq.failOnErr(err, "failed to connect rabbitMq")

	// 获取 channel
	rabbitMq.channel, err = rabbitMq.conn.Channel()
	rabbitMq.failOnErr(err, "failed to open a channel")

	return rabbitMq
}

// 路由模式发送消息
func (r *RabbitMQ)PublishRouting(message string) {
	// 1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// direct
		"direct",
		true,
		false,
		false,
		false,
		nil,
		)
	r.failOnErr(err, "Failed to declare an exchange")

	// 2. 发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.key,
		false,
		false,
		amqp.Publishing{
			ContentType:"text/plain",
			Body: []byte(message),
		})
}

// 路由模式接受消息
func (r *RabbitMQ) ReceiveRouting() {
	// 1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
		)
	r.failOnErr(err, "Failed to declare an exchange")

	// 2. 尝试创建队列， 这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		// 随机生产队列名称
		"",
		false,
		false,
		true,
		false,
		nil,
		)
	r.failOnErr(err, "Failed to declare a queue")

	// 绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		r.key,
		r.Exchange,
		false,
		nil)

	// 消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
		)

	forever := make(chan bool)

	go func() {
		for d:= range messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("退出请按 CTRL+C\n")
	<- forever
}


// 话题模式 创建RabbitMQ 实例
func NewRabbitMQTopic(exchangeName string, routingKey string)*RabbitMQ {
	// 创建RabbitMQ 实例
	rabbitMq := NewRabbitMQ("", exchangeName, routingKey)
	var err  error

	// 获取 connection
	rabbitMq.conn, err = amqp.Dial(rabbitMq.MqUrl)
	rabbitMq.failOnErr(err, "failed to connection rabbitMq")

	// 获取 channel
	rabbitMq.channel, err = rabbitMq.conn.Channel()
	rabbitMq.failOnErr(err, "failed to open a channel")

	return rabbitMq
}

// 话题模式发送消息
func (r *RabbitMQ)PublishTopic(message string) {
	// 1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
		)
	r.failOnErr(err, "failed to declare an exchange")

	// 2. 发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(message),
		})
}

// 话题模式接受消息
// 要注意 key 规则：其中“*”用于匹配一个单词，“#”用于匹配多个单词（可以是零个）
//匹配 news.* 表示匹配 news.stock, news.stock.000001需要用news.#才能匹配到

func (r *RabbitMQ) ReceiveTopic() {
	// 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
		)
	r.failOnErr(err, "failed to declare an exchange")

	// 尝试创建队列，这里注意队列名称不用写
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
		)
	r.failOnErr(err, "failed to declare a queue")

	// 绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		r.key,
		r.Exchange,
		false,
		nil,
		)

	// 消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
		)

	forever := make(chan bool)

	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("退出请按 CTRL+C\n")

	<- forever
}

package main

import (
	"awesomeProject/common/mqUtils"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	initConsumer1()
}

const (
	NORMAL_EXCHANGE   = "normal_exchange"
	DEAD_EXCHANGE     = "dead_exchange"
	NORMAL_QUEUE      = "normal_queue"
	DEAD_QUEUE        = "dead_queue"
	NORMAL_ROUTINGKEY = "normal"
	DEAD_ROUTINGKEY   = "dead"
	EXCHANGE_TYPE     = "direct"
)

func initConsumer1() {
	url := fmt.Sprintf("amqp://%s:%s@%s:5672/", mqUtils.MQ_USER, mqUtils.MQ_PWD, mqUtils.MQ_ADDR)
	con, err := amqp.Dial(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	ch, err := con.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 1. 声明 normal_exchange
	err = ch.ExchangeDeclare(
		NORMAL_EXCHANGE,
		EXCHANGE_TYPE,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. 声明 dead_exchange
	err = ch.ExchangeDeclare(
		DEAD_EXCHANGE,
		EXCHANGE_TYPE,
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 3. 声明 normal_queue
	args := amqp.Table{
		"x-dead-letter-exchange":    DEAD_EXCHANGE,   //死信队列交换机
		"x-dead-letter-routing-key": DEAD_ROUTINGKEY, //死信队列routing key
		//"x-max-length":              2,               //测试二：队列最大长度
	}
	normal_q, err := ch.QueueDeclare(
		NORMAL_QUEUE,
		false,
		false,
		false,
		false,
		args,
	)

	// 4. 声明 dead_queue
	dead_q, err := ch.QueueDeclare(
		DEAD_QUEUE,
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//5. 队列绑定
	err = ch.QueueBind(
		normal_q.Name,
		NORMAL_ROUTINGKEY,
		NORMAL_EXCHANGE,
		false,
		nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ch.QueueBind(
		dead_q.Name,
		DEAD_ROUTINGKEY,
		DEAD_EXCHANGE,
		false,
		nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//6. 接受normal_queue中的消息
	msgs, err := ch.Consume(
		normal_q.Name,
		"c1",
		false, //测试三：拒绝特定消息
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	forever := make(chan interface{})
	go func() {
		for msg := range msgs {
			if string(msg.Body) == "msg-5" {
				fmt.Println("c1 abandon normal msg -->", string(msg.Body))
				msg.Nack(false, false) //拒绝消息
			} else {
				fmt.Println("c1 receive normal msg -->", string(msg.Body))
				msg.Ack(false) //确认消息
			}

		}
	}()

	<-forever
}

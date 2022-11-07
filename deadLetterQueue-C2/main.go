package main

import (
	"awesomeProject/common/mqUtils"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	DEAD_QUEUE = "dead_queue"
)

func main() {
	initConsumer2()
}

func initConsumer2() {
	//获取mq连接
	url := fmt.Sprintf("amqp://%s:%s@%s:5672/", mqUtils.MQ_USER, mqUtils.MQ_PWD, mqUtils.MQ_ADDR)
	con, err := amqp.Dial(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	//获取channel
	ch, err := con.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}

	//消费死信队列中的消息
	msgs, err := ch.Consume(
		DEAD_QUEUE,
		"c2",
		true,
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
			fmt.Println("c2 receive dead msg -->", string(msg.Body))
		}
	}()

	<-forever
}

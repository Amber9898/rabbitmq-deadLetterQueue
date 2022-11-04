package main

import (
	"awesomeProject/common/mqUtils"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

const (
	NORMAL_EXCHANGE   = "normal_exchange"
	DEAD_EXCHANGE     = "dead_exchange"
	NORMAL_QUEUE      = "normal_queue"
	DEAD_QUEUE        = "dead_queue"
	NORMAL_ROUTINGKEY = "normal"
	DEAD_ROUTINGKEY   = "dead"
	EXCHANGE_TYPE     = "direct"
)

var mqCh *amqp.Channel

func main() {
	initPublisher()
	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("msg-%d", i)
		publishMessage(msg)
	}

}

func publishMessage(msg string) {
	if mqCh == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mqCh.PublishWithContext(
		ctx,
		NORMAL_EXCHANGE,
		NORMAL_ROUTINGKEY,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
}
func initPublisher() {
	url := fmt.Sprintf("amqp://%s:%s@%s:5672/", mqUtils.MQ_USER, mqUtils.MQ_PWD, mqUtils.MQ_ADDR)
	con, err := amqp.Dial(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	mqCh, err = con.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Copyright 2021 Zubin. All rights reserved.

package main

import (
	"fmt"
	"log"

	"github.com/zubinzhang/taskqueue"
)

func main() {
	mq, err := taskqueue.NewRabbitMQ("amqp://admin:password@localhost:5672/", taskqueue.Option{
		QueueName: "test",
	})
	if err != nil {
		fmt.Printf("%+v", err)
	}

	defer mq.Destory()

	onConsumed := func(msg taskqueue.Message) error {
		log.Printf("Received a message: %s", msg.Payload)
		return nil
	}

	err = mq.Consume(onConsumed)

	if err != nil {
		fmt.Printf("%+v", err)
	}
}

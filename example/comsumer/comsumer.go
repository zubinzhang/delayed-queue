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
	}).Connect()
	if err != nil {
		fmt.Printf("%+v", err)
	}

	defer mq.Destroy()

	onConsumed := func(msg taskqueue.Message) error {
		log.Printf("Received a message: %s", msg.Payload)

		return nil
	}

	err = mq.Consume(onConsumed)

	if err != nil {
		fmt.Printf("%+v", err)
	}
}

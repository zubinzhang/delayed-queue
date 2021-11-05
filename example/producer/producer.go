// Copyright 2021 Zubin. All rights reserved.

package main

import (
	"fmt"
	"time"

	queue "github.com/zubinzhang/delayed-queue"
)

func main() {
	publisher, err := queue.NewPublisher(
		"amqp://admin:password@localhost:5672/",
		queue.WithPublisherOptionsExchange("test"),
		queue.WithPublisherOptionsQueue("test"),
		queue.WithPublisherOptionsKey("test"),
	)
	if err != nil {
		fmt.Printf("%+v", err)
	}

	defer publisher.Disconnect()

	err = publisher.Publish([]byte("Hello"), 3*time.Second)
	if err != nil {
		fmt.Printf("%+v", err)
	}
}

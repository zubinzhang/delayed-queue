// Copyright 2021 Zubin. All rights reserved.

package main

import (
	"fmt"
	"log"

	queue "github.com/zubinzhang/delayed-queue"
)

func main() {
	comsumer, err := queue.NewComsumer(
		"amqp://admin:password@localhost:5672/",
		queue.WithComsumerOptionsExchange("test"),
		queue.WithComsumerOptionsQueue("test"),
		queue.WithComsumerOptionsKey("test"),
	)
	if err != nil {
		fmt.Printf("%+v", err)
	}

	defer comsumer.Disconnect()

	deliveries := comsumer.Consume()

	for d := range deliveries {
		log.Printf("Received a message: %s", d.Body)
	}
}

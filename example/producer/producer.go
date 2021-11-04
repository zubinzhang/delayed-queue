// Copyright 2021 Zubin. All rights reserved.

package main

import (
	"fmt"
	"time"

	"github.com/zubinzhang/taskqueue"
)

func main() {
	publisher, err := taskqueue.NewPublisher("amqp://admin:password@localhost:5672/", taskqueue.WithPublisherOptionsSererviceName("test"))
	if err != nil {
		fmt.Printf("%+v", err)
	}

	defer publisher.Disconnect()

	err = publisher.Publish("test", []byte("Hello"), 3*time.Second)
	if err != nil {
		fmt.Printf("%+v", err)
	}
}

// Copyright 2021 Zubin. All rights reserved.

package main

import (
	"fmt"
	"time"

	"github.com/zubinzhang/delayedqueue"
)

func main() {
	publisher, err := delayedqueue.NewPublisher("amqp://admin:password@localhost:5672/", delayedqueue.WithPublisherOptionsSererviceName("test"))
	if err != nil {
		fmt.Printf("%+v", err)
	}

	defer publisher.Disconnect()

	err = publisher.Publish("test", []byte("Hello"), 3*time.Second)
	if err != nil {
		fmt.Printf("%+v", err)
	}
}

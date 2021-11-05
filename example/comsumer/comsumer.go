// Copyright 2021 Zubin. All rights reserved.

package main

import (
	"fmt"
	"log"

	"github.com/zubinzhang/delayedqueue"
)

func main() {
	comsumer, err := delayedqueue.NewComsumer("amqp://admin:password@localhost:5672/", delayedqueue.WithComsumerOptionsSererviceName("test"))
	if err != nil {
		fmt.Printf("%+v", err)
	}

	defer comsumer.Disconnect()

	onConsumed := func(msg delayedqueue.Message) error {
		log.Printf("Received a message: %s", msg.Payload)

		return nil
	}

	err = comsumer.Consume(onConsumed)

	if err != nil {
		fmt.Printf("%+v", err)
	}
}

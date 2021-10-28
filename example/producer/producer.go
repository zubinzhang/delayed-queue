// Copyright 2021 Zubin. All rights reserved.

package main

import (
	"fmt"

	"github.com/zubinzhang/taskqueue"
)

func main() {
	mq, err := taskqueue.NewTaskQueue("amqp://admin:password@localhost:5672/", taskqueue.Option{
		ServiceName: "test",
	}).Connect()
	if err != nil {
		fmt.Printf("%+v", err)
	}

	defer mq.Destroy()

	err = mq.Publish("test", "Hello", 0)
	if err != nil {
		fmt.Printf("%+v", err)
	}
}

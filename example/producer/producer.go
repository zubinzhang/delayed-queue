// Copyright 2021 Zubin. All rights reserved.

package main

import (
	"fmt"
	"time"

	"github.com/zubinzhang/taskqueue"
)

func main() {
	mq, err := taskqueue.New("amqp://admin:password@localhost:5672/", taskqueue.SererviceName("test")).Connect()
	if err != nil {
		fmt.Printf("%+v", err)
	}

	fmt.Println(int64(time.Second / time.Millisecond))

	defer mq.Destroy()

	err = mq.Publish("test", "Hello", 10*time.Second)
	if err != nil {
		fmt.Printf("%+v", err)
	}
}

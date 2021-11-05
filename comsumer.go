// Copyright 2021 Zubin. All rights reserved.

package delayedq

import (
	"time"

	"github.com/marmotedu/log"
	"github.com/streadway/amqp"
)

type Comsumer struct {
	DelayedQueue
}

// New Publisher returns a new publisher with an open channel.
func NewComsumer(url string, options ...ComsumerOptions) (*Comsumer, error) {
	var err error

	c := &Comsumer{
		DelayedQueue: getDefaultTaskQueue(),
	}
	for _, option := range options {
		option(c)
	}

	c.url = url

	err = c.connect()
	if err != nil {
		return nil, err
	}

	go c.handleDisconnect()

	err = c.init()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// consume from the queue.
func (c *Comsumer) Consume() <-chan amqp.Delivery {
	deliveries := make(chan amqp.Delivery)

	go func() {
		for {
			msgs, err := c.channel.Consume(
				c.queue, // queue
				"",      // consumer
				true,    // auto-ack
				false,   // exclusive
				false,   // no-local
				false,   // no-wait
				nil,     // args
			)
			if err != nil {
				log.Errorf("consume failed, err: %v", err)
				time.Sleep(3 * time.Second)
				continue
			}

			for msg := range msgs {
				deliveries <- msg
			}

			// sleep before IsClose call. closed flag may not set before sleep.
			time.Sleep(time.Second)

			if _, ok := <-c.quitChan; ok {
				break
			}
		}
	}()

	return deliveries
}

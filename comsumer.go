// Copyright 2021 Zubin. All rights reserved.

package delayedqueue

import (
	"time"

	"github.com/marmotedu/log"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/zubinzhang/delayedqueue/protos"
	"google.golang.org/protobuf/proto"
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

func (c *Comsumer) announceQueue() (<-chan amqp.Delivery, error) {
	deliveries, err := c.channel.Consume(
		c.workQueue, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to register a consumer")
	}
	return deliveries, nil
}

// consume from the queue.
func (c *Comsumer) Consume(handler func(msg Message) error) (err error) {

	forever := make(chan bool)

	go func() {
		for {
			deliveries, err := c.announceQueue()
			if err != nil {
				log.Errorf("consume failed, err: %v", err)
				time.Sleep(3 * time.Second)
				continue
			}

			for d := range deliveries {
				payload := &protos.Payload{}
				err = proto.Unmarshal(d.Body, payload)
				if err != nil {
					log.Fatalf("Unmarshaling error: ", err)
				}
				handler(Message{
					CorrelationId: d.CorrelationId,
					MessageId:     d.MessageId,
					Timestamp:     d.Timestamp,
					JobName:       payload.JobName,
					Payload:       payload.Body,
					Priority:      int(*payload.Priority),
				})
			}
		}
	}()

	log.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}

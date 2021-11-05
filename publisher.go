// Copyright 2021 Zubin. All rights reserved.

package delayedq

import (
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Publisher struct {
	DelayedQueue
	retry    int
	priority int
}

// New Publisher returns a new publisher with an open channel.
func NewPublisher(url string, options ...PublisherOptions) (*Publisher, error) {
	p := &Publisher{
		DelayedQueue: getDefaultTaskQueue(),
		retry:        0,
		priority:     0,
	}
	for _, option := range options {
		option(p)
	}

	p.url = url

	err := p.connect()
	if err != nil {
		return nil, err
	}

	go p.handleDisconnect()

	err = p.init()
	if err != nil {
		return nil, err
	}

	return p, nil
}

// publish the provided data over the connection
func (p *Publisher) Publish(body []byte, delay time.Duration) (err error) {
	headers := make(amqp.Table)
	if delay != 0 {
		headers["x-delay"] = int64(delay / time.Millisecond)
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         body,
		Headers:      headers,
	}

	err = p.channel.Publish(
		p.exchange, // exchange
		p.key,      // routing key
		false,      // mandatory
		false,      // immediate
		msg,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to publish a message")
	}

	return nil
}

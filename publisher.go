// Copyright 2021 Zubin. All rights reserved.

package taskqueue

import (
	"time"

	"github.com/go-basic/uuid"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/zubinzhang/taskqueue/protos"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Publisher struct {
	TaskQueue
	retry    int
	priority int
}

// New Publisher returns a new publisher with an open channel.
func NewPublisher(url string, options ...PublisherOptions) (*Publisher, error) {
	p := &Publisher{
		TaskQueue: getDefaultTaskQueue(),
		retry:     0,
		priority:  0,
	}
	for _, option := range options {
		option(p)
	}

	p.url = url

	err := p.connect()
	if err != nil {
		return nil, err
	}

	err = p.init()
	if err != nil {
		return nil, err
	}

	return p, nil
}

// publish the provided data over the connection
func (p *Publisher) Publish(jobName string, body []byte, delay time.Duration) (err error) {
	payload, err := proto.Marshal(&protos.Payload{
		Id:        uuid.New(),
		JobName:   jobName,
		Timestamp: timestamppb.Now(),
		Body:      body,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to marshaling payload")
	}

	headers := make(amqp.Table)
	if delay != 0 {
		headers["x-delay"] = int64(delay / time.Millisecond)
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         payload,
		Headers:      headers,
	}

	err = p.channel.Publish(
		p.exchange,  // exchange
		p.workQueue, // routing key
		false,       // mandatory
		false,       // immediate
		msg,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to publish a message")
	}

	return nil
}

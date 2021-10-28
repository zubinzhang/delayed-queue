// Copyright 2021 Zubin. All rights reserved.

package taskqueue

import (
	"fmt"
	"log"
	"time"

	"github.com/go-basic/uuid"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/zubinzhang/taskqueue/protos"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultPrefetchCount = 1
	defaultMessageTTL    = int(10 * time.Second)
)

type Message struct {
	Payload       []byte
	CorrelationId string
	MessageId     string
	Timestamp     time.Time
}

type Option struct {
	QueueName     string
	Exchange      string
	Key           string
	serviceName   string
	messageTTL    int
	prefetchCount int
}

type TaskQueue struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
	option  Option
}

// new rabbitmq instance.
func NewRabbitMQ(amqpUrl string, config Option) *TaskQueue {
	if config.QueueName == "" {
		panic("QueueName can not be empty")
	}
	return &TaskQueue{url: amqpUrl, option: config}
}

// connect connection and channel
func (q *TaskQueue) Connect() (*TaskQueue, error) {
	var err error

	q.conn, err = amqp.Dial(q.url)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to RabbitMQ")
	}

	q.channel, err = q.conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open a channel")
	}
	return q, nil
}

// close channel and connection.
func (q *TaskQueue) Destroy() {
	q.channel.Close()
	q.conn.Close()
}

// publish a msg.
func (q *TaskQueue) Publish(jobName, body string) error {
	_, err := q.channel.QueueDeclare(
		q.option.QueueName, // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return errors.Wrap(err, "Failed to declare a queue")
	}

	payload, err := proto.Marshal(&protos.Payload{
		Id:        uuid.New(),
		Timestamp: timestamppb.Now(),
		JobName:   jobName,
		Body:      []byte(body),
	})
	if err != nil {
		errors.Wrap(err, "Failed to marshaling payload")
	}

	err = q.channel.Publish(
		q.option.Exchange,  // exchange
		q.option.QueueName, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        payload,
		})
	if err != nil {
		return errors.Wrap(err, "Failed to publish a message")
	}

	return nil
}

// consume msg.
func (q *TaskQueue) Consume(handler func(msg Message) error) error {
	_, err := q.channel.QueueDeclare(
		q.option.QueueName, // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return errors.Wrap(err, "Failed to declare a queue")
	}

	msgs, err := q.channel.Consume(
		q.option.QueueName, // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		return errors.Wrap(err, "Failed to register a consumer")
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			payload := &protos.Payload{}
			err = proto.Unmarshal(d.Body, payload)
			if err != nil {
				log.Fatal("unmarshaling error: ", err)
			}
			fmt.Println(payload)
			handler(Message{CorrelationId: d.CorrelationId, Payload: payload.Body, MessageId: d.MessageId, Timestamp: d.Timestamp})
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}

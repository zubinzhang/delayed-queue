// Copyright 2021 Zubin. All rights reserved.

package taskqueue

import (
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
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
	PrefetchCount int
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	option  Option
	url     string
}

// new rabbitmq instance
func NewRabbitMQ(amqpUrl string, config Option) (*RabbitMQ, error) {
	if config.QueueName == "" {
		panic("QueueName can not be empty")
	}
	rabbitmq := &RabbitMQ{url: amqpUrl, option: config}
	var err error

	rabbitmq.conn, err = amqp.Dial(rabbitmq.url)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to RabbitMQ")
	}

	rabbitmq.channel, err = rabbitmq.conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open a channel")
	}

	return rabbitmq, nil
}

// close channel and connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// publish a msg
func (r *RabbitMQ) Publish(message string) error {
	_, err := r.channel.QueueDeclare(
		r.option.QueueName, // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return errors.Wrap(err, "Failed to declare a queue")
	}

	err = r.channel.Publish(
		r.option.Exchange,  // exchange
		r.option.QueueName, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return errors.Wrap(err, "Failed to publish a message")
	}

	return nil
}

// consume msg
func (r *RabbitMQ) Consume(handler func(msg Message) error) error {
	q, err := r.channel.QueueDeclare(
		r.option.QueueName, // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return errors.Wrap(err, "Failed to declare a queue")
	}

	msgs, err := r.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return errors.Wrap(err, "Failed to register a consumer")
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			handler(Message{CorrelationId: d.CorrelationId, Payload: d.Body, MessageId: d.MessageId, Timestamp: d.Timestamp})
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

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
	defaultPrefetchCount  = 1
	defaultMessageTTL     = int(10 * time.Second)
	DELAYED_EXCHANGE_TYPE = "x-delayed-message"
	DELAYED_TYPE          = "direct"
)

type Message struct {
	Payload       []byte
	CorrelationId string
	MessageId     string
	Timestamp     time.Time
}

type Option struct {
	ServiceName   string
	MessageTTL    int
	PrefetchCount int
}

type TaskQueue struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	url         string
	option      Option
	workQueue   string
	failedQueue string
	exchange    string
	key         string
}

// new rabbitmq instance.
func NewTaskQueue(amqpUrl string, config Option) *TaskQueue {
	if config.ServiceName == "" {
		panic("serviceName can not be empty")
	}
	return &TaskQueue{
		url:         amqpUrl,
		option:      config,
		exchange:    fmt.Sprintf("%s_exchange", config.ServiceName),
		workQueue:   fmt.Sprintf("%s_work_queue", config.ServiceName),
		failedQueue: fmt.Sprintf("%s_failed_queue", config.ServiceName),
	}
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
func (q *TaskQueue) Publish(jobName, body string, delay int) error {
	args := make(amqp.Table)
	args["x-delayed-type"] = DELAYED_TYPE
	err := q.channel.ExchangeDeclare(
		q.exchange,            // name
		DELAYED_EXCHANGE_TYPE, // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		args,                  // arguments
	)

	if err != nil {
		return errors.Wrap(err, "Failed to declare a exchange")
	}
	// _, err = q.channel.QueueDeclare(
	// 	q.option.queueName, // name
	// 	false,              // durable
	// 	false,              // delete when unused
	// 	false,              // exclusive
	// 	false,              // no-wait
	// 	nil,                // arguments
	// )
	// if err != nil {
	// 	return errors.Wrap(err, "Failed to declare a queue")
	// }

	payload, err := proto.Marshal(&protos.Payload{
		Id:        uuid.New(),
		Timestamp: timestamppb.Now(),
		JobName:   jobName,
		Body:      []byte(body),
	})
	if err != nil {
		errors.Wrap(err, "Failed to marshaling payload")
	}

	headers := make(amqp.Table)
	if delay != 0 {
		headers["x-delay"] = delay
	}

	err = q.channel.Publish(
		q.exchange,  // exchange
		q.workQueue, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			ContentType:  "text/plain",
			Body:         payload,
			Headers:      headers,
		})
	if err != nil {
		return errors.Wrap(err, "Failed to publish a message")
	}

	return nil
}

// consume msg.
func (q *TaskQueue) Consume(handler func(msg Message) error) error {
	args := make(amqp.Table)
	args["x-delayed-type"] = DELAYED_TYPE
	err := q.channel.ExchangeDeclare(
		q.exchange,            // name
		DELAYED_EXCHANGE_TYPE, // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		args,                  // arguments
	)
	if err != nil {
		return errors.Wrap(err, "Failed to declare a exchange")
	}

	_, err = q.channel.QueueDeclare(
		q.workQueue, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return errors.Wrap(err, "Failed to declare a queue")
	}

	err = q.channel.QueueBind(q.workQueue, q.workQueue, q.exchange, false, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to bind queue")
	}

	msgs, err := q.channel.Consume(
		q.workQueue, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
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

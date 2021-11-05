// Copyright 2021 Zubin. All rights reserved.

package delayedqueue

import (
	"fmt"
	"time"

	"github.com/marmotedu/log"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const (
	DEFAULT_PREFETCH_COUNT = 1
	DEFAULT_SERVICE_NAME   = "task_queue"
	DELAYED_EXCHANGE_TYPE  = "x-delayed-message"
	DELAYED_TYPE           = "direct"
)

type Message struct {
	JobName       string
	Payload       []byte
	CorrelationId string
	MessageId     string
	Timestamp     time.Time
	Priority      int
}

type DelayedQueue struct {
	connection    *amqp.Connection
	channel       *amqp.Channel
	url           string
	exchange      string
	workQueue     string
	failedQueue   string
	prefetchCount int
	closeChan     chan *amqp.Error
	quitChan      chan bool
}

func getDefaultTaskQueue() DelayedQueue {
	return DelayedQueue{
		exchange:      fmt.Sprintf("%s_exchange", DEFAULT_SERVICE_NAME),
		workQueue:     fmt.Sprintf("%s_work_queue", DEFAULT_SERVICE_NAME),
		failedQueue:   fmt.Sprintf("%s_failed_queue", DEFAULT_SERVICE_NAME),
		prefetchCount: DEFAULT_PREFETCH_COUNT,
		quitChan:      make(chan bool),
	}
}

func (dq *DelayedQueue) connect() error {
	var err error

	dq.connection, err = amqp.Dial(dq.url)
	if err != nil {
		return errors.Wrap(err, "Failed to connect to RabbitMQ")
	}

	dq.channel, err = dq.connection.Channel()
	if err != nil {
		return errors.Wrap(err, "Failed to open a channel")
	}

	log.Debug("Connect to rabbitMQ established")

	dq.closeChan = make(chan *amqp.Error)
	dq.connection.NotifyClose(dq.closeChan)

	return nil
}

// handleDisconnect handle a disconnection trying to reconnect every 5s.
func (dq *DelayedQueue) handleDisconnect() {
	for {
		select {
		case errChan := <-dq.closeChan:
			if errChan != nil {
				log.Errorf("RabbitMQ disconnection: %v", errChan)
			}
		case <-dq.quitChan:
			dq.connection.Close()
			dq.channel.Close()
			dq.connection = nil
			dq.channel = nil
			log.Debug("RabbitMQ has been shut down...")
			dq.quitChan <- true
			return
		}

		log.Info("Trying to reconnect to rabbitMQ...")
		time.Sleep(5 * time.Second)

		if err := dq.connect(); err != nil {
			log.Errorf("RabbitMQ connect error: %v", err)
		}

		log.Info("Connect to rabbitMQ established")

		if err := dq.init(); err != nil {
			log.Errorf("RabbitMQ init error: %v", err)
		}
	}
}

// declare exchange and queue if not exist
func (dq *DelayedQueue) init() (err error) {
	args := make(amqp.Table)
	args["x-delayed-type"] = DELAYED_TYPE
	err = dq.channel.ExchangeDeclare(
		dq.exchange,           // name
		DELAYED_EXCHANGE_TYPE, // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // noWait
		args,                  // arguments
	)
	if err != nil {
		return errors.Wrap(err, "Failed to declare a exchange")
	}

	_, err = dq.channel.QueueDeclare(
		dq.workQueue, // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return errors.Wrap(err, "Failed to declare a queue")
	}

	err = dq.channel.QueueBind(dq.workQueue, dq.workQueue, dq.exchange, false, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to bind queue")
	}

	err = dq.channel.Qos(
		dq.prefetchCount, // prefetch count
		0,                // prefetch size
		false,            // global
	)
	if err != nil {
		return errors.Wrap(err, "Failed to set QoS")
	}

	return nil
}

// Disconnect the channel and connection
func (dq *DelayedQueue) Disconnect() {
	dq.quitChan <- true
	log.Debug("shutting down rabbitMQ's connection...")
	<-dq.quitChan
}

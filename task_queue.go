// Copyright 2021 Zubin. All rights reserved.

package taskqueue

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
}

type TaskQueue struct {
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

func getDefaultTaskQueue() TaskQueue {
	return TaskQueue{
		exchange:      fmt.Sprintf("%s_exchange", DEFAULT_SERVICE_NAME),
		workQueue:     fmt.Sprintf("%s_work_queue", DEFAULT_SERVICE_NAME),
		failedQueue:   fmt.Sprintf("%s_failed_queue", DEFAULT_SERVICE_NAME),
		prefetchCount: DEFAULT_PREFETCH_COUNT,
		quitChan:      make(chan bool),
	}
}

func (tq *TaskQueue) connect() error {
	var err error

	tq.connection, err = amqp.Dial(tq.url)
	if err != nil {
		return errors.Wrap(err, "Failed to connect to RabbitMQ")
	}

	tq.channel, err = tq.connection.Channel()
	if err != nil {
		return errors.Wrap(err, "Failed to open a channel")
	}

	log.Info("Connect to rabbitMQ established")

	tq.closeChan = make(chan *amqp.Error)
	tq.connection.NotifyClose(tq.closeChan)

	return nil
}

// handleDisconnect handle a disconnection trying to reconnect every 5s.
func (tq *TaskQueue) handleDisconnect() {
	for {
		select {
		case errChan := <-tq.closeChan:
			if errChan != nil {
				log.Errorf("RabbitMQ disconnection: %v", errChan)
			}
		case <-tq.quitChan:
			tq.connection.Close()
			tq.channel.Close()
			tq.connection = nil
			tq.channel = nil
			log.Debug("RabbitMQ has been shut down...")
			tq.quitChan <- true
			return
		}

		log.Info("Trying to reconnect to rabbitMQ...")
		time.Sleep(5 * time.Second)

		if err := tq.connect(); err != nil {
			log.Errorf("RabbitMQ connect error: %v", err)
		}

		if err := tq.init(); err != nil {
			log.Errorf("RabbitMQ init error: %v", err)
		}
	}
}

// declare exchange and queue if not exist
func (tq *TaskQueue) init() (err error) {
	args := make(amqp.Table)
	args["x-delayed-type"] = DELAYED_TYPE
	err = tq.channel.ExchangeDeclare(
		tq.exchange,
		DELAYED_EXCHANGE_TYPE,
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to declare a exchange")
	}

	_, err = tq.channel.QueueDeclare(
		tq.workQueue, // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return errors.Wrap(err, "Failed to declare a queue")
	}

	err = tq.channel.QueueBind(tq.workQueue, tq.workQueue, tq.exchange, false, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to bind queue")
	}

	err = tq.channel.Qos(
		tq.prefetchCount, // prefetch count
		0,                // prefetch size
		false,            // global
	)
	if err != nil {
		return errors.Wrap(err, "Failed to set QoS")
	}

	return nil
}

// Disconnect the channel and connection
func (tq *TaskQueue) Disconnect() {
	tq.quitChan <- true
	log.Debug("shutting down rabbitMQ's connection...")
	<-tq.quitChan
}

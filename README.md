# delayed-queue-go-sdk

A delayed queue base [rabbitmq](https://www.rabbitmq.com/), use the [rabbitmq_delayed_message_exchange](https://github.com/rabbitmq/rabbitmq-delayed-message-exchange)'s plugin.

We also handle a disconnection and reconnection of rabbitmq trying every 5 sec if a new connection is available.

## Installation

Start up rabbitmq in local.

```bash
docker compose up -d
```

## Usage

### Publish Message

You can publish the message by calling `publish()`

```go
	publisher, err := queue.NewPublisher(
		"amqp://admin:password@localhost:5672/",
		queue.WithPublisherOptionsExchange("test"),
		queue.WithPublisherOptionsQueue("test"),
		queue.WithPublisherOptionsKey("test"),
	)
	if err != nil {
		fmt.Printf("%+v", err)
	}

	defer publisher.Disconnect()

	err = publisher.Publish([]byte("Hello"), 3*time.Second)
	if err != nil {
		fmt.Printf("%+v", err)
	}
```

## Subscribe Message

```go
comsumer, err := queue.NewComsumer(
		"amqp://admin:password@localhost:5672/",
		queue.WithComsumerOptionsExchange("test"),
		queue.WithComsumerOptionsQueue("test"),
		queue.WithComsumerOptionsKey("test"),
	)
	if err != nil {
		fmt.Printf("%+v", err)
	}

	defer comsumer.Disconnect()

	deliveries := comsumer.Consume()

	for d := range deliveries {
		log.Printf("Received a message: %s", d.Body)
	}
```

## Questions & Suggestions

Please open an issue [here](https://github.com/zubinzhang/delayed-queue-go-sdk/issues).

## License

[MIT](https://github.com/zubinzhang/delayed-queue-go-sdk/blob/main/LICENSE)

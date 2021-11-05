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
publisher, err := taskqueue.NewPublisher("amqp://127.0.0.1:5672", taskqueue.WithPublisherOptionsSererviceName("test"))
	if err != nil {
		fmt.Printf("%+v", err)
	}

	defer publisher.Disconnect()

	err = publisher.Publish("test", []byte("Hello"), 3*time.Second)
	if err != nil {
		fmt.Printf("%+v", err)
	}
```

## Subscribe Message

```go
comsumer, err := taskqueue.NewComsumer("amqp://admin:password@localhost:5672/", taskqueue.WithComsumerOptionsSererviceName("test"))
	if err != nil {
		fmt.Printf("%+v", err)
	}

	defer comsumer.Disconnect()

	onConsumed := func(msg taskqueue.Message) error {
		log.Printf("Received a message: %s", msg.Payload)

		return nil
	}

	err = comsumer.Consume(onConsumed)

	if err != nil {
		fmt.Printf("%+v", err)
	}
```

## Questions & Suggestions

Please open an issue [here](https://github.com/zubinzhang/delayed-queue-go-sdk/issues).

## License

[MIT](https://github.com/zubinzhang/delayed-queue-go-sdk/blob/main/LICENSE)

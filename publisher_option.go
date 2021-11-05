// Copyright 2021 Zubin. All rights reserved.

package delayedq

type PublisherOptions func(*Publisher)

// WithPublisherOptionsPriority returns a function that sets the exchange.
func WithPublisherOptionsExchange(exchange string) PublisherOptions {
	return func(p *Publisher) {
		p.DelayedQueue.exchange = exchange
	}
}

// WithPublisherOptionsPriority returns a function that sets the queue.
func WithPublisherOptionsQueue(queue string) PublisherOptions {
	return func(p *Publisher) {
		p.DelayedQueue.queue = queue
	}
}

// WithPublisherOptionsPriority returns a function that sets the binding key.
func WithPublisherOptionsKey(key string) PublisherOptions {
	return func(p *Publisher) {
		p.DelayedQueue.key = key
	}
}

// WithPublisherOptionsPriority returns a function that sets the priority.
func WithPublisherOptionsPriority(priority int) PublisherOptions {
	return func(p *Publisher) {
		p.priority = priority
	}
}

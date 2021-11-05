// Copyright 2021 Zubin. All rights reserved.

package delayedq

type ComsumerOptions func(*Comsumer)

// WithComsumerOptionsPriority returns a function that sets the exchange.
func WithComsumerOptionsExchange(exchange string) ComsumerOptions {
	return func(c *Comsumer) {
		c.DelayedQueue.exchange = exchange
	}
}

// WithComsumerOptionsPriority returns a function that sets the queue.
func WithComsumerOptionsQueue(queue string) ComsumerOptions {
	return func(c *Comsumer) {
		c.DelayedQueue.queue = queue
	}
}

// WithComsumerOptionsPriority returns a function that sets the binding key.
func WithComsumerOptionsKey(key string) ComsumerOptions {
	return func(c *Comsumer) {
		c.DelayedQueue.key = key
	}
}

func WithComsumerOptionsPrefetchCount(prefetchCount int) ComsumerOptions {
	return func(c *Comsumer) {
		c.DelayedQueue.prefetchCount = prefetchCount
	}
}

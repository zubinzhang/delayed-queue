// Copyright 2021 Zubin. All rights reserved.

package delayedqueue

import (
	"fmt"
)

type ComsumerOptions func(*Comsumer)

func WithComsumerOptionsSererviceName(name string) ComsumerOptions {
	return func(c *Comsumer) {
		c.DelayedQueue.exchange = fmt.Sprintf("%s_exchange", name)
		c.DelayedQueue.workQueue = fmt.Sprintf("%s_work_queue", name)
		c.DelayedQueue.failedQueue = fmt.Sprintf("%s_failed_queue", name)
	}
}

func WithComsumerOptionsPrefetchCount(prefetchCount int) ComsumerOptions {
	return func(c *Comsumer) {
		c.DelayedQueue.prefetchCount = prefetchCount
	}
}

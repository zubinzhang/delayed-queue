// Copyright 2021 Zubin. All rights reserved.

package taskqueue

import (
	"fmt"
)

type ComsumerOptions func(*Comsumer)

func WithComsumerOptionsSererviceName(name string) ComsumerOptions {
	return func(c *Comsumer) {
		c.TaskQueue.exchange = fmt.Sprintf("%s_exchange", name)
		c.TaskQueue.workQueue = fmt.Sprintf("%s_work_queue", name)
		c.TaskQueue.failedQueue = fmt.Sprintf("%s_failed_queue", name)
	}
}

func WithComsumerOptionsPrefetchCount(prefetchCount int) ComsumerOptions {
	return func(c *Comsumer) {
		c.TaskQueue.prefetchCount = prefetchCount
	}
}

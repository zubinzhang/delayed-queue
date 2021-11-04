// Copyright 2021 Zubin. All rights reserved.

package taskqueue

import (
	"fmt"
	"math"
)

type PublisherOptions func(*Publisher)

const (
	MINIMUM_RETRY = 0
	MAXIMUM_RETRY = 10
)

func WithPublisherOptionsSererviceName(name string) PublisherOptions {
	return func(p *Publisher) {
		p.TaskQueue.exchange = fmt.Sprintf("%s_exchange", name)
		p.TaskQueue.workQueue = fmt.Sprintf("%s_work_queue", name)
		p.TaskQueue.failedQueue = fmt.Sprintf("%s_failed_queue", name)
	}
}

// WithPublisherOptionsRetry returns a function that sets the retry count.
// maximum retry count is 10.
func WithPublisherOptionsRetry(retry int) PublisherOptions {
	return func(p *Publisher) {
		p.retry = int(math.Min(math.Max(MINIMUM_RETRY, float64(retry)), MAXIMUM_RETRY))
	}
}

// WithPublisherOptionsPriority returns a function that sets the priority.
func WithPublisherOptionsPriority(priority int) PublisherOptions {
	return func(p *Publisher) {
		p.priority = priority
	}
}

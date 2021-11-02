// Copyright 2021 Zubin. All rights reserved.

package taskqueue

import (
	"fmt"
)

type Option func(*TaskQueue)

func SererviceName(name string) Option {
	return func(tq *TaskQueue) {
		tq.exchange = fmt.Sprintf("%s_exchange", name)
		tq.workQueue = fmt.Sprintf("%s_work_queue", name)
		tq.failedQueue = fmt.Sprintf("%s_failed_queue", name)
	}
}

func PrefetchCount(count int) Option {
	return func(tq *TaskQueue) {
		tq.prefetchCount = count
	}
}

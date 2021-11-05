// Copyright 2021 Zubin. All rights reserved.

package delayedqueue

import (
	"fmt"
)

type Options func(*DelayedQueue)

func SererviceName(name string) Options {
	return func(tq *DelayedQueue) {
		tq.exchange = fmt.Sprintf("%s_exchange", name)
		tq.workQueue = fmt.Sprintf("%s_work_queue", name)
		tq.failedQueue = fmt.Sprintf("%s_failed_queue", name)
	}
}

func PrefetchCount(count int) Options {
	return func(tq *DelayedQueue) {
		tq.prefetchCount = count
	}
}

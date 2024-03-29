/*
 * @Author: cnzf1
 * @Date: 2023-01-04 09:46:56
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-28 19:35:53
 * @Description:
 */
package queue

import (
	"container/heap"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cnzf1/gocore/lang"
)

// DelayQueue is an unbounded blocking queue of *Delayed* elements, in which
// an element can only be taken when its delay has expired. The head of the
// queue is the *Delayed* element whose delay expired furthest in the past.
type DelayQueue struct {
	C chan lang.AnyType

	mu sync.Mutex
	pq PriorityQueue

	// Similar to the sleeping state of runtime.timers.
	sleeping int32
	wakeupC  chan lang.PlaceholderType
}

// NewDelayQueue creates an instance of delayQueue with the specified size.
func NewDelayQueue(size int) *DelayQueue {
	return &DelayQueue{
		C:       make(chan lang.AnyType),
		pq:      NewPriorityQueue(size),
		wakeupC: make(chan lang.PlaceholderType),
	}
}

// Offer inserts the element into the current queue.
func (dq *DelayQueue) Offer(elem lang.AnyType, expiration int64) {
	item := &PriorityQueueItem{Value: elem, Priority: expiration}

	dq.mu.Lock()
	heap.Push(&dq.pq, item)
	index := item.Index
	dq.mu.Unlock()

	if index == 0 {
		// A new item with the earliest expiration is added.
		if atomic.CompareAndSwapInt32(&dq.sleeping, 1, 0) {
			dq.wakeupC <- struct{}{}
		}
	}
}

// Poll starts an infinite loop, in which it continually waits for an element
// to expire and then send the expired element to the channel C.
func (dq *DelayQueue) Poll(exitC chan lang.PlaceholderType, nowF func() int64) {
	for {
		now := nowF()

		dq.mu.Lock()
		item, delta := dq.pq.PeekAndShift(now)
		if item == nil {
			// No items left or at least one item is pending.

			// We must ensure the atomicity of the whole operation, which is
			// composed of the above PeekAndShift and the following StoreInt32,
			// to avoid possible race conditions between Offer and Poll.
			atomic.StoreInt32(&dq.sleeping, 1)
		}
		dq.mu.Unlock()

		if item == nil {
			if delta == 0 {
				// No items left.
				select {
				case <-dq.wakeupC:
					// Wait until a new item is added.
					continue
				case <-exitC:
					goto exit
				}
			} else if delta > 0 {
				// At least one item is pending.
				select {
				case <-dq.wakeupC:
					// A new item with an "earlier" expiration than the current "earliest" one is added.
					continue
				case <-time.After(time.Duration(delta) * time.Millisecond):
					// The current "earliest" item expires.

					// Reset the sleeping state since there's no need to receive from wakeupC.
					if atomic.SwapInt32(&dq.sleeping, 0) == 0 {
						// A caller of Offer() is being blocked on sending to wakeupC,
						// drain wakeupC to unblock the caller.
						<-dq.wakeupC
					}
					continue
				case <-exitC:
					goto exit
				}
			}
		}

		select {
		case dq.C <- item.Value:
			// The expired element has been sent out successfully.
		case <-exitC:
			goto exit
		}
	}

exit:
	// Reset the states
	atomic.StoreInt32(&dq.sleeping, 0)
}

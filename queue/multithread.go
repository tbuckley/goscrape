package queue

import (
	"net/url"
	"sync"
)

// SimpleQueue is a basic implementation of a queue that uses a linked list.
// It is thread-safe.
type MultithreadQueue struct {
	queue Queue
	lock  *sync.Mutex
	cond  *sync.Cond
}

func NewMultithreadQueue(queue Queue) *MultithreadQueue {
	q := new(MultithreadQueue)
	q.queue = queue
	q.lock = new(sync.Mutex)
	q.cond = sync.NewCond(q.lock)
	return q
}

// Push adds an element to the back of the queue.
func (q *MultithreadQueue) Push(page *url.URL) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	err := q.queue.Push(page)
	if err != nil {
		return err
	}
	q.cond.Signal()
	return nil
}

// Pop returns the first element in the queue. If the queue is empty, the call
// blocks until it has something to return.
func (q *MultithreadQueue) Pop() (*url.URL, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.queue.Pop()
}

// PopBlock returns the first element in the queue. If the queue is empty, the
// call blocks until it has something to return.
func (q *MultithreadQueue) PopBlock() *url.URL {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Wait for the queue to be populated
	for q.queue.Length() == 0 {
		q.cond.Wait()
	}

	page, _ := q.queue.Pop()
	return page
}

// Length returns the number of elements in the queue.
func (q *MultithreadQueue) Length() int {
	return q.queue.Length()
}

type MultithreadPrioQueue struct {
	*MultithreadQueue
	prioQueue PriorityQueue
}

func NewMultithreadPrioQueue(prioQueue PriorityQueue) *MultithreadPrioQueue {
	q := new(MultithreadPrioQueue)
	q.MultithreadQueue = NewMultithreadQueue(prioQueue)
	q.prioQueue = prioQueue
	return q
}

func (q *MultithreadPrioQueue) PushPriority(page *url.URL, prio uint) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	err := q.prioQueue.PushPriority(page, prio)
	if err != nil {
		return err
	}
	q.cond.Signal()
	return nil
}

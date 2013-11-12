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

func NewMultithreadQueue(name string) Queue {
	q := new(MultithreadQueue)
	q.queue = NewSimpleQueue(name)
	q.lock = new(sync.Mutex)
	q.cond = sync.NewCond(q.lock)
	return q
}

// Push adds an element to the back of the queue.
func (q *MultithreadQueue) Push(page *url.URL) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.queue.Push(page)
	q.cond.Signal()
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

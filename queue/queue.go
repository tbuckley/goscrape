package queue

import (
	"net/url"
)

// Queue objects are FIFO data structures.
type Queue interface {
	// Push adds an element to the back of the queue.
	Push(page *url.URL) error

	// Pop returns the first element in the queue. If the queue is empty, the
	// second argument will be false.
	Pop() (*url.URL, bool)

	// Length returns the number of elements in the queue.
	Length() int
}

type PriorityQueue interface {
	Queue

	// PushPriority adds an element to the queue. Elements with higher priorities
	// will be removed before those with lower ones.
	PushPriority(page *url.URL, prio uint) error
}

type AsyncQueue interface {
	Queue

	// PopBlock returns the first element in the queue. If the queue is empty,
	// the call blocks until it has something to return.
	PopBlock() *url.URL
}

type AsyncPrioQueue interface {
	PriorityQueue

	PopBlock() *url.URL
}

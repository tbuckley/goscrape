package queue

import (
	"net/url"
)

// ComboQueue lets you combine multiple queues into one, pulling from them in
// the order they are added.
type ComboQueue struct {
	queues []Queue
}

// NewComboQueue creates a ComboQueue from the queues provided. It accesses
// them in the order they are given.
func NewComboQueue(queues ...Queue) Queue {
	return &ComboQueue{queues}
}

// Push is not supported on ComboQueue. It is provided to meet the interface.
func (q *ComboQueue) Push(page *url.URL) {
	panic("Cannot push to ComboQueue")
}

// Pop returns an element from the first non-empty queue it finds.
func (q *ComboQueue) Pop() (*url.URL, bool) {
	for _, queue := range q.queues {
		if queue.Length() > 0 {
			return queue.Pop()
		}
	}
	return nil, false
}

// Length returns the total length of all the combined queues.
func (q *ComboQueue) Length() int {
	sum := 0
	for _, queue := range q.queues {
		sum += queue.Length()
	}
	return sum
}

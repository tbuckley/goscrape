package queue

import (
	"expvar"
	"net/url"
)

type node struct {
	Next *node
	Page *url.URL
}

// SimpleQueue is a basic implementation of a queue that uses a linked list.
type SimpleQueue struct {
	head    *node
	tail    *node
	size    int
	varSize *expvar.Int
}

func NewSimpleQueue(name string) Queue {
	s := new(SimpleQueue)
	s.varSize = expvar.NewInt(name)
	return s
}

// Push adds an element to the back of the queue.
func (q *SimpleQueue) Push(page *url.URL) {
	n := &node{Page: page}
	if q.head == nil {
		q.head = n
	} else {
		q.tail.Next = n
	}
	q.tail = n
	q.size += 1
	q.varSize.Add(1)
}

// Pop returns the first element in the queue. If the queue is empty, the
// second argument will be false.
func (q *SimpleQueue) Pop() (*url.URL, bool) {
	if q.size == 0 {
		return nil, false
	}

	q.size -= 1
	q.varSize.Add(-1)
	if q.head == q.tail {
		n := q.head
		q.head, q.tail = nil, nil
		return n.Page, true
	} else {
		n := q.head
		q.head = q.head.Next
		return n.Page, true
	}
}

// Length returns the number of elements in the queue.
func (q *SimpleQueue) Length() int {
	return q.size
}

package queue

import (
	"expvar"
	"fmt"
	"net/url"
)

// PrioQueue lets you combine multiple queues into one, pulling from them in
// the order they are added.
type PrioQueue struct {
	queues  []Queue
	bitmap  Bitmap
	length  int
	varSize *expvar.Int
}

// NewPrioQueue creates a PrioQueue from the queues provided. It accesses
// them in the order they are given.
func NewPrioQueue(name string) *PrioQueue {
	q := new(PrioQueue)
	q.queues = make([]Queue, 32)
	for i := 0; i < 32; i++ {
		q.queues[i] = NewSimpleQueue(fmt.Sprintf("%s_%02d", name, i))
	}
	q.varSize = expvar.NewInt(name)
	return q
}

func (q *PrioQueue) Push(page *url.URL) error {
	return q.PushPriority(page, 0)
}

func (q *PrioQueue) PushPriority(page *url.URL, prio uint) error {
	q.bitmap.Set(prio, true)
	q.queues[prio].Push(page)
	q.length += 1
	q.varSize.Add(1)
	return nil
}

// Pop returns an element from the highest-priority non-empty queue.
func (q *PrioQueue) Pop() (*url.URL, bool) {
	hiprio, err := q.bitmap.HighBit()
	if err != nil {
		return nil, false
	}

	page, ok := q.queues[hiprio].Pop()
	q.bitmap.Set(hiprio, q.queues[hiprio].Length() != 0)
	q.length--
	q.varSize.Add(-1)
	return page, ok
}

// Length returns the total length of all the combined queues.
func (q *PrioQueue) Length() int {
	return q.length
}

// // //

type BitmapEmptyError uint32

func (e BitmapEmptyError) Error() string {
	return "Bitmap is empty"
}

type Bitmap struct {
	n uint32
}

var MultiplyDeBruijnBitPosition2 = []uint{0, 1, 28, 2, 29, 14, 24, 3, 30, 22,
	20, 15, 25, 17, 4, 8, 31, 27, 13, 23, 21, 19, 16, 7, 26, 12, 18, 6, 11, 5, 10, 9}

func (b *Bitmap) HighBit() (uint, error) {
	x := b.n
	x |= (x >> 1)
	x |= (x >> 2)
	x |= (x >> 4)
	x |= (x >> 8)
	x |= (x >> 16)
	x = x - (x >> 1)

	if x == 0 {
		return 0, BitmapEmptyError(b.n)
	}

	return MultiplyDeBruijnBitPosition2[uint32(x*0x077CB531)>>27], nil
}

func (b *Bitmap) Set(bit uint, value bool) {
	if value {
		b.n |= 1 << bit
	} else {
		b.n &= ^(1 << bit)
	}
}

package utils

import (
	"container/list"
	"sync"
)

// Queue represents the FIFO queue.
type Queue struct {
	dead bool
	l    *list.List
}

// A thread safe version of Queue.
type AsyncQueue struct {
	*Queue
	lock sync.RWMutex
}

// Returns an initialized Queue.
func NewQueue() *Queue {
	return &Queue{l: list.New()}
}

// Returns an initialized SyncQueue
func NewAsyncQueue() *AsyncQueue {
	return &AsyncQueue{Queue: NewQueue()}
}

func (q *Queue) Kill() {
	q.dead = true
}

func (q *Queue) Dead() bool {
	return q.dead
}

// Pushes a new item to the back of the Queue.
func (q *Queue) Push(o interface{}) {
	q.l.PushBack(o)
}

// Removes an item from the front of the Queue and returns it's value or nil.
func (q *Queue) Pop() interface{} {
	e := q.l.Front()
	if e == nil {
		return nil
	}

	return q.l.Remove(e)
}

// Checks to see if the Queue is empty.
func (q *Queue) IsEmpty() bool {
	return q.l.Len() == 0
}

// Returns the current length of the Queue.
func (q *Queue) Len() int {
	return q.l.Len()
}

// Returns the item at the front of the Queue or nil.
// The item is a *list.Element from the 'container/list' package.
func (q *Queue) Front() *list.Element {
	return q.l.Front()
}

// Returns the item after e or nil it is the last item or nil.
// The item is a *list.Element from the 'container/list' package.
// Even though it is possible to call e.Next() directly, don't. This behavior
// may not be supported moving forward.
func (q *Queue) Next(e *list.Element) *list.Element {
	if e == nil {
		return e
	}

	return e.Next()
}

// Same as Push for Queue, except it is thread safe.
func (q *AsyncQueue) Push(o interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.Queue.Push(o)
}

// Same as Pop for Queue, except it is thread safe.
func (q *AsyncQueue) Pop() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.Queue.Pop()
}

// Same as IsEmpty for Queue, except it is thread safe.
func (q *AsyncQueue) IsEmpty() bool {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return q.Queue.IsEmpty()
}

// Same as Len for Queue, except it is thread safe.
func (q *AsyncQueue) Len() int {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return q.Queue.Len()
}

// Same as Front for Queue, except it is thread safe.
func (q *AsyncQueue) Front() *list.Element {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return q.Queue.Front()
}

// Same as Next for Queue, except it is thread safe.
func (q *AsyncQueue) Next(e *list.Element) *list.Element {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return q.Queue.Next(e)
}

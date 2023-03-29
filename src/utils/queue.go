package utils

import (
	"container/heap"
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

// An Item is something we manage in a priority queue.
type Item struct {
	value    interface{} // The value of the item; arbitrary.
	priority int64       // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A myHeap implements heap.Interface and holds Items.
type myHeap []*Item

func (mh myHeap) Len() int { return len(mh) }

func (mh myHeap) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return mh[i].priority > mh[j].priority
}

func (mh myHeap) Swap(i, j int) {
	mh[i], mh[j] = mh[j], mh[i]
	mh[i].index = i
	mh[j].index = j
}

func (mh *myHeap) Push(x interface{}) {
	n := len(*mh)
	item := x.(*Item)
	item.index = n
	*mh = append(*mh, item)
}

func (mh *myHeap) Pop() interface{} {
	old := *mh
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*mh = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (mh *myHeap) Update(item *Item, value string, priority int64) {
	item.value = value
	item.priority = priority
	heap.Fix(mh, item.index)
}

type PriorityQueue struct {
	h *myHeap
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{h: &myHeap{}}
}

func (this *PriorityQueue) PushWithPriority(value interface{}, priority int64) {
	item := &Item{value: value, priority: priority}
	heap.Push(this.h, item)
}

func (this *PriorityQueue) Push(value interface{}) {
	item := &Item{value: value, priority: -1}
	heap.Push(this.h, item)
}

func (this *PriorityQueue) Pop() interface{} {
	if this.h.Len() < 1 {
		return nil
	}
	return heap.Pop(this.h).(*Item).value
}

type AsyncPriorityQueue struct {
	h    *myHeap
	lock sync.RWMutex
}

func NewAsyncPriorityQueue() *AsyncPriorityQueue {
	return &AsyncPriorityQueue{h: &myHeap{}}
}

func (this *AsyncPriorityQueue) PushWithPriority(value interface{}, priority int64) {
	this.lock.Lock()
	defer this.lock.Unlock()
	item := &Item{value: value, priority: priority}
	heap.Push(this.h, item)
}

func (this *AsyncPriorityQueue) Push(value interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	item := &Item{value: value, priority: -1}
	heap.Push(this.h, item)
}

func (this *AsyncPriorityQueue) Pop() interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()
	if this.h.Len() < 1 {
		return nil
	}
	return heap.Pop(this.h).(*Item).value
}

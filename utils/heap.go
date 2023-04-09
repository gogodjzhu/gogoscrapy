package utils

// An Item is something we manage in a priority queue.
type Item struct {
	value    interface{} // The value of the item; arbitrary.
	priority int64       // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

func NewItem(value interface{}, priority int64) *Item {
	return &Item{
		value:    value,
		priority: priority,
	}
}

func (i *Item) GetValue() interface{} {
	return i.value
}

func (i *Item) SetValue(value interface{}) *Item {
	i.value = value
	return i
}

// A Heap implements heap.Interface and holds Items.
type Heap []*Item

func NewHeap() *Heap {
	h := make(Heap, 0)
	return &h
}

func (mh Heap) First() *Item { return mh[0] }

func (mh Heap) Len() int { return len(mh) }

func (mh Heap) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return mh[i].priority > mh[j].priority
}

func (mh Heap) Swap(i, j int) {
	mh[i], mh[j] = mh[j], mh[i]
	mh[i].index = i
	mh[j].index = j
}

func (mh *Heap) Push(x interface{}) {
	n := len(*mh)
	item := x.(*Item)
	item.index = n
	*mh = append(*mh, item)
}

func (mh *Heap) Pop() interface{} {
	old := *mh
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*mh = old[0 : n-1]
	return item
}

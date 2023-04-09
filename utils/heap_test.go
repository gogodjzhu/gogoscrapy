package utils

import (
	"container/heap"
	"testing"
)

func TestHeap(t *testing.T) {
	h := NewHeap()
	for i := 0; i < 1000; i++ {
		heap.Push(h, NewItem(i, int64(999-i)))
	}
	for i := 0; i < 1000; i++ {
		item := heap.Pop(h).(*Item)
		if item.GetValue().(int) != 999-i {
			t.Error("heap pop error")
			t.FailNow()
		}
	}
}

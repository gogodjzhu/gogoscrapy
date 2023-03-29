package utils

import (
	"container/heap"
	"sync"
	"testing"
)

func TestHeap(t *testing.T) {
	h := myHeap{}
	for i := 0; i < 7; i++ {
		h = append(h, &Item{value: i + 1, priority: int64(i + 1), index: i})
	}

	heap.Init(&h)
	i := heap.Pop(&h).(*Item).priority
	if i != 7 {
		t.Error("Heap pop the element which is not the max")
	}
	heap.Push(&h, &Item{value: 4, priority: 4})
	i = heap.Pop(&h).(*Item).priority
	if i != 6 {
		t.Error("Heap pop the element which is not the max")
	}
	heap.Push(&h, &Item{value: 11, priority: 11})
	i = heap.Pop(&h).(*Item).priority
	if i != 11 {
		t.Error("Heap pop the element which is not the max")
	}
}

func TestPriorityQueue(t *testing.T) {
	priorityQueue := NewPriorityQueue()
	priorityQueue.PushWithPriority("1", 1)
	priorityQueue.PushWithPriority("3", 3)
	priorityQueue.PushWithPriority("2", 2)
	if priorityQueue.Pop().(string) != "3" {
		t.Error("test failed @ PriorityQueue")
	}
	priorityQueue.PushWithPriority("4", 4)
	if priorityQueue.Pop().(string) != "4" {
		t.Error("test failed @ PriorityQueue")
	}
}

func TestAsyncPriorityQueue(t *testing.T) {
	wg := sync.WaitGroup{}
	asyncPriorityQueue := NewAsyncPriorityQueue()
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				asyncPriorityQueue.PushWithPriority(j, int64(j))
			}
		}()
	}
	wg.Wait()
	for i := 19; i >= 0; i-- {
		for j := 0; j < 20; j++ {
			if asyncPriorityQueue.Pop().(int) != i {
				t.Error("test failed @ AsyncPriorityQueue")
			}
		}
	}
	if asyncPriorityQueue.Pop() != nil {
		t.Error("test failed @ AsyncPriorityQueue")
	}

}

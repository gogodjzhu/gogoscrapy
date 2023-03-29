package utils

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

func TestConcurrentSet_Add(t *testing.T) {
	testRound := func() {
		m := NewAsyncSet()
		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(index int) {
				for j := 0; j < 100; j++ {
					str := strconv.Itoa(index*10000 + j)
					m.Add(str)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
		if m.Size() != 10000 {
			t.Error("test failed @ TestConcurrentSet_Add")
		}
	}
	for r := 0; r < 50; r++ {
		testRound()
	}
}

func TestConcurrentSet_Remove(t *testing.T) {
	testRound := func() {
		m := NewAsyncSet()
		for i := 0; i < 10000; i++ {
			m.Add(i)
		}

		var optCnt int32
		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(index int) {
				for step := 0; step < 100; step++ {
					m.Remove(index*100 + step)
					atomic.AddInt32(&optCnt, 1)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
		if !m.IsEmpty() {
			t.Error("test failed @ TestConcurrentSet_Remove")
		}
		if optCnt != 10000 || m.Size() != 0 {
			t.Error("test failed @ TestConcurrentSet_Remove")
		}
	}
	for r := 0; r < 50; r++ {
		testRound()
	}
}

func TestConcurrentSet_Clear(t *testing.T) {
	testRound := func() {
		var cnt int
		m := NewAsyncSet()
		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(index int) {
				for j := 0; j < 100; j++ {
					str := strconv.Itoa(index*10000 + j)
					m.Add(str)
				}
				wg.Done()
			}(i)
			if i == 10 || i == 40 || i == 80 {
				cnt += m.Clear()
			}
		}
		wg.Wait()
		cnt += m.Size()
		if cnt != 10000 {
			t.Error("test failed @ TestConcurrentSet_Clear")
		}
	}
	for r := 0; r < 50; r++ {
		testRound()
	}
}

package utils

import "sync"

type Set struct {
	m map[interface{}]bool
}

func NewSet() *Set {
	return &Set{m: make(map[interface{}]bool)}
}
func (this *Set) Size() int {
	return len(this.m)
}

func (this *Set) IsEmpty() bool {
	return len(this.m) < 1
}

func (this *Set) Contains(e interface{}) bool {
	_, ok := this.m[e]
	return ok
}

func (this *Set) Add(e interface{}) bool {
	_, ok := this.m[e]
	if !ok {
		this.m[e] = true
		return true
	}
	return false
}

func (this *Set) Remove(e interface{}) bool {
	_, ok := this.m[e]
	if ok {
		delete(this.m, e)
		return true
	}
	return false
}

func (this *Set) Clear() int {
	cnt := len(this.m)
	this.m = make(map[interface{}]bool)
	return cnt
}

type AsyncSet struct {
	sync.RWMutex
	m map[interface{}]bool
}

func NewAsyncSet() *AsyncSet {
	return &AsyncSet{m: make(map[interface{}]bool)}
}

func (rm *AsyncSet) Size() int {
	rm.RLock()
	defer rm.RUnlock()
	return len(rm.m)
}

func (rm *AsyncSet) IsEmpty() bool {
	rm.RLock()
	defer rm.RUnlock()
	return len(rm.m) < 1
}

func (rm *AsyncSet) Contains(e interface{}) bool {
	rm.RLock()
	defer rm.RUnlock()
	_, ok := rm.m[e]
	return ok
}

//add specified element to set, return true if this set did not contain this element.
func (rm *AsyncSet) Add(e interface{}) bool {
	rm.Lock()
	defer rm.Unlock()
	_, ok := rm.m[e]
	if !ok {
		rm.m[e] = true
		return true
	}
	return false
}

func (rm *AsyncSet) Remove(e interface{}) bool {
	rm.Lock()
	defer rm.Unlock()
	_, ok := rm.m[e]
	if ok {
		delete(rm.m, e)
		return true
	}
	return false
}

func (rm *AsyncSet) Clear() int {
	rm.Lock()
	defer rm.Unlock()
	cnt := len(rm.m)
	rm.m = make(map[interface{}]bool)
	return cnt
}

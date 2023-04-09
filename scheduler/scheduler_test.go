package scheduler

import (
	"context"
	"github.com/gogodjzhu/gogoscrapy/entity"
	"sync"
	"testing"
	"time"
)

func TestMemScheduler_Basic(t *testing.T) {
	ctx, wg := context.Background(), &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)
	var scheduler = NewMemScheduler(ctx, wg)
	for i := 0; i < 100; i++ {
		go func(i int) {
			scheduler.PushChan() <- &entity.Request{Url: "abc", Priority: int64(i)}
		}(i)
	}
	pullChan := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			_, ok := <-scheduler.PullChan()
			pullChan <- ok
		}()
	}
	pullCnt := 0
	for {
		select {
		case <-time.After(5 * time.Second):
			t.Error("pull timeout")
			t.FailNow()
		case <-pullChan:
			pullCnt++
			if pullCnt == 100 {
				cancel()
				wg.Wait()
				return
			}
		}
	}
}

func TestMemScheduler_Priority(t *testing.T) {
	ctx, wg := context.Background(), &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)
	var scheduler = NewMemScheduler(ctx, wg)

	pushWg := &sync.WaitGroup{}
	pushWg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			scheduler.PushChan() <- &entity.Request{Url: "abc", Priority: int64(99 - i)}
			pushWg.Done()
		}(i)
	}
	pushWg.Wait() // all request pushed, and order by priority

	pullCnt := 0
	for {
		select {
		case <-time.After(5 * time.Second):
			t.Error("pull timeout")
			t.FailNow()
		case req := <-scheduler.PullChan():
			pullCnt++
			if req.GetPriority() != int64(100-pullCnt) { // TODO still some error
				t.Errorf("priority error: expect %d, got %d", 100-pullCnt, req.GetPriority())
				t.Fail()
			}
			if pullCnt == 100 {
				cancel()
				wg.Wait()
				return
			}
		}
	}
}

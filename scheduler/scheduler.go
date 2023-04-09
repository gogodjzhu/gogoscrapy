package scheduler

import (
	"container/heap"
	"context"
	ent "github.com/gogodjzhu/gogoscrapy/entity"
	"github.com/gogodjzhu/gogoscrapy/utils"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type IScheduler interface {
	PushChan() chan ent.IRequest
	PullChan() chan ent.IRequest
	CacheSize() int
}

type MemScheduler struct {
	remover        DuplicateRemover
	inChan         chan ent.IRequest
	outChan        chan ent.IRequest
	bufferRequests *utils.Heap
}

func NewMemScheduler(ctx context.Context, wg *sync.WaitGroup) *MemScheduler {
	scheduler := &MemScheduler{
		remover:        NewMemDuplicateRemover(),
		inChan:         make(chan ent.IRequest),
		outChan:        make(chan ent.IRequest),
		bufferRequests: utils.NewHeap(),
	}
	lock := &sync.Mutex{}
	// a helper function to pull request from bufferRequests
	pullBuffer := func() ent.IRequest {
		lock.Lock()
		defer lock.Unlock()
		if scheduler.bufferRequests.Len() == 0 {
			return nil
		}
		next := heap.Pop(scheduler.bufferRequests).(*utils.Item)
		req := next.GetValue().(ent.IRequest)
		if isDup, err := scheduler.remover.IsDuplicate(req); err != nil {
			log.Warnf("failed to check duplicate err: %+v", err)
			return req
		} else if !isDup {
			return req
		}
		return nil
	}

	// a helper function to push request into bufferRequests
	pushBuffer := func(req ent.IRequest) {
		lock.Lock()
		defer lock.Unlock()
		item := utils.NewItem(req, req.GetPriority())
		heap.Push(scheduler.bufferRequests, item)
	}
	go func(ctx context.Context, wg *sync.WaitGroup) { // out from scheduler
		wg.Add(1)
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Infof("scheduler@out is closed by context, drop remain buffer size: %d", scheduler.bufferRequests.Len())
				return
			default:
				next := pullBuffer()
				if next == nil {
					time.Sleep(1 * time.Millisecond)
					continue
				}
				select {
				case <-ctx.Done():
					log.Infof("scheduler@out is closed by context, drop remain buffer size: %d", scheduler.bufferRequests.Len())
					return
				case scheduler.outChan <- next:
					// inChan is unbuffered, this will block until received
				}
			}
		}
	}(ctx, wg)
	go func(ctx context.Context, wg *sync.WaitGroup) { // into scheduler
		wg.Add(1)
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				log.Infof("scheduler@in is closed by context")
				return
			case req := <-scheduler.inChan:
				pushBuffer(req)
			}
		}
	}(ctx, wg)
	return scheduler
}

func (qs *MemScheduler) PushChan() chan ent.IRequest {
	return qs.inChan
}

func (qs *MemScheduler) PullChan() chan ent.IRequest {
	return qs.outChan
}

func (qs *MemScheduler) CacheSize() int {
	return qs.bufferRequests.Len()
}

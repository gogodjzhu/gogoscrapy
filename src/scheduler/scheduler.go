package scheduler

import (
	"gogoscrapy/src/entity"
	"gogoscrapy/src/utils"
	"net/http"
	"sunteng/commons/log"
	"sync/atomic"
	"time"
)

type IScheduler interface {
	entity.Closeable
	Push(request entity.IRequest)
	Poll() entity.IRequest
	Size() int
}

const (
	StatRunning = iota
	StatClosing
	StatClosed
)

type QueueScheduler struct {
	remover DuplicateRemover
	queue   *utils.AsyncQueue
	stat    atomic.Value
}

func NewQueueScheduler() *QueueScheduler {
	running := atomic.Value{}
	running.Store(StatRunning)
	return &QueueScheduler{
		stat:    running,
		queue:   utils.NewAsyncQueue(),
		remover: NewMemDuplicateRemover(),
	}
}

func (this *QueueScheduler) Push(req entity.IRequest) {
	if this.stat.Load() != StatRunning {
		return
	}
	if noNeedToRemoveDuplicate(req) || !this.remover.IsDuplicate(req) {
		log.Logf("push req, %+s", req.GetUrl())
		this.queue.Push(req)
	} else if req.IsRetry() {
		log.Logf("push retry req, %+s", req.GetUrl())
		this.queue.Push(req)
	}
}

func (this *QueueScheduler) Poll() entity.IRequest {
	ret := this.queue.Pop()
	if ret != nil {
		req := ret.(entity.IRequest)
		log.Logf("poll req, %+s", req.GetUrl())
		return req
	} else {
		return nil
	}
}

func (this *QueueScheduler) Size() int {
	return this.queue.Len()
}

func (this *QueueScheduler) Close() error {
	this.stat.Store(StatClosing)
	for !this.queue.IsEmpty() {
		time.Sleep(1 * time.Second)
		log.Logf("schedule waiting queue clear, left size:%d", this.queue.Len())
	}
	this.stat.Store(StatClosed)
	return nil
}

func (this *QueueScheduler) IsClose() bool {
	return this.stat.Load() == StatClosed
}

func noNeedToRemoveDuplicate(request entity.IRequest) bool {
	return http.MethodPost == request.GetMethod()
}

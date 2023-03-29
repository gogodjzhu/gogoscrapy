package scheduler

import (
	entity2 "github.com/gogodjzhu/gogoscrapy/entity"
	utils2 "github.com/gogodjzhu/gogoscrapy/utils"
	"net/http"
	"sync/atomic"
	"time"
)

var LOG = utils2.NewLogger()

type IScheduler interface {
	entity2.Closeable
	Push(request entity2.IRequest)
	Poll() entity2.IRequest
	Size() int
}

const (
	StatRunning = iota
	StatClosing
	StatClosed
)

type QueueScheduler struct {
	remover            DuplicateRemover
	queue              *utils2.AsyncQueue
	asyncPriorityQueue *utils2.AsyncPriorityQueue
	stat               atomic.Value
}

func NewQueueScheduler() *QueueScheduler {
	running := atomic.Value{}
	running.Store(StatRunning)
	return &QueueScheduler{
		stat:               running,
		queue:              utils2.NewAsyncQueue(),
		asyncPriorityQueue: utils2.NewAsyncPriorityQueue(),
		remover:            NewMemDuplicateRemover(),
	}
}

func (this *QueueScheduler) Push(req entity2.IRequest) {
	if this.stat.Load() != StatRunning {
		return
	}
	if noNeedToRemoveDuplicate(req) || !this.remover.IsDuplicate(req) {
		LOG.Infof("push req, %+s", req.GetUrl())
		if req.GetPriority() > 0 {
			this.pushWithPriority(req, req.GetPriority())
		} else {
			this.queue.Push(req)
		}
	} else if req.IsRetry() {
		LOG.Infof("push retry req, %+s", req.GetUrl())
		if req.GetPriority() > 0 {
			this.pushWithPriority(req, req.GetPriority())
		} else {
			this.queue.Push(req)
		}
	}
}

func (this *QueueScheduler) pushWithPriority(req entity2.IRequest, priority int64) {
	this.asyncPriorityQueue.PushWithPriority(req, priority)
}

func (this *QueueScheduler) Poll() entity2.IRequest {
	ret := this.asyncPriorityQueue.Pop()
	if ret != nil {
		req := ret.(entity2.IRequest)
		LOG.Infof("poll req, %+s", req.GetUrl())
		return req
	}
	ret = this.queue.Pop()
	if ret != nil {
		req := ret.(entity2.IRequest)
		LOG.Infof("poll req, %+s", req.GetUrl())
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
		LOG.Infof("schedule waiting queue clear, left size:%d", this.queue.Len())
	}
	this.stat.Store(StatClosed)
	return nil
}

func (this *QueueScheduler) IsClose() bool {
	return this.stat.Load() == StatClosed
}

func noNeedToRemoveDuplicate(request entity2.IRequest) bool {
	return http.MethodPost == request.GetMethod()
}

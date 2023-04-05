package scheduler

import (
	ent "github.com/gogodjzhu/gogoscrapy/entity"
	u "github.com/gogodjzhu/gogoscrapy/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync/atomic"
	"time"
)

type IScheduler interface {
	ent.Closeable
	Push(request ent.IRequest)
	Poll() ent.IRequest
	Size() int
}

const (
	StatRunning = iota
	StatClosing
	StatClosed
)

type QueueScheduler struct {
	remover            DuplicateRemover
	queue              *u.AsyncQueue
	asyncPriorityQueue *u.AsyncPriorityQueue
	stat               atomic.Value
}

func NewQueueScheduler() *QueueScheduler {
	running := atomic.Value{}
	running.Store(StatRunning)
	return &QueueScheduler{
		stat:               running,
		queue:              u.NewAsyncQueue(),
		asyncPriorityQueue: u.NewAsyncPriorityQueue(),
		remover:            NewMemDuplicateRemover(),
	}
}

func (this *QueueScheduler) Push(req ent.IRequest) {
	if this.stat.Load() != StatRunning {
		return
	}
	if noNeedToRemoveDuplicate(req) || !this.remover.IsDuplicate(req) {
		log.Debugf("push req, %+s", req.GetUrl())
		if req.GetPriority() > 0 {
			this.pushWithPriority(req, req.GetPriority())
		} else {
			this.queue.Push(req)
		}
	} else if req.IsRetry() {
		log.Debugf("push retry req, %+s", req.GetUrl())
		if req.GetPriority() > 0 {
			this.pushWithPriority(req, req.GetPriority())
		} else {
			this.queue.Push(req)
		}
	}
}

func (this *QueueScheduler) pushWithPriority(req ent.IRequest, priority int64) {
	this.asyncPriorityQueue.PushWithPriority(req, priority)
}

func (this *QueueScheduler) Poll() ent.IRequest {
	ret := this.asyncPriorityQueue.Pop()
	if ret != nil {
		req := ret.(ent.IRequest)
		log.Infof("poll req, %+s", req.GetUrl())
		return req
	}
	ret = this.queue.Pop()
	if ret != nil {
		req := ret.(ent.IRequest)
		log.Infof("poll req, %+s", req.GetUrl())
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
		log.Infof("schedule waiting queue clear, left size:%d", this.queue.Len())
	}
	this.stat.Store(StatClosed)
	return nil
}

func (this *QueueScheduler) IsClose() bool {
	return this.stat.Load() == StatClosed
}

func noNeedToRemoveDuplicate(request ent.IRequest) bool {
	return http.MethodPost == request.GetMethod()
}

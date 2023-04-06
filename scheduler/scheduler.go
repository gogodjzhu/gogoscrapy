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

func (qs *QueueScheduler) Push(req ent.IRequest) {
	if qs.stat.Load() != StatRunning {
		return
	}
	isDuplicate, _ := qs.remover.IsDuplicate(req)
	if noNeedToRemoveDuplicate(req) || !isDuplicate {
		log.Debugf("push req, %+s", req.GetUrl())
		if req.GetPriority() > 0 {
			qs.pushWithPriority(req, req.GetPriority())
		} else {
			qs.queue.Push(req)
		}
	} else if req.IsRetry() {
		log.Debugf("push retry req, %+s", req.GetUrl())
		if req.GetPriority() > 0 {
			qs.pushWithPriority(req, req.GetPriority())
		} else {
			qs.queue.Push(req)
		}
	}
}

func (qs *QueueScheduler) pushWithPriority(req ent.IRequest, priority int64) {
	qs.asyncPriorityQueue.PushWithPriority(req, priority)
}

func (qs *QueueScheduler) Poll() ent.IRequest {
	ret := qs.asyncPriorityQueue.Pop()
	if ret != nil {
		req := ret.(ent.IRequest)
		log.Infof("poll req, %+s", req.GetUrl())
		return req
	}
	ret = qs.queue.Pop()
	if ret != nil {
		req := ret.(ent.IRequest)
		log.Infof("poll req, %+s", req.GetUrl())
		return req
	} else {
		return nil
	}
}

func (qs *QueueScheduler) Size() int {
	return qs.queue.Len()
}

func (qs *QueueScheduler) Close() error {
	qs.stat.Store(StatClosing)
	for !qs.queue.IsEmpty() {
		time.Sleep(1 * time.Second)
		log.Infof("schedule waiting queue clear, left size:%d", qs.queue.Len())
	}
	qs.stat.Store(StatClosed)
	return nil
}

func (qs *QueueScheduler) IsClose() bool {
	return qs.stat.Load() == StatClosed
}

func noNeedToRemoveDuplicate(request ent.IRequest) bool {
	return http.MethodPost == request.GetMethod()
}

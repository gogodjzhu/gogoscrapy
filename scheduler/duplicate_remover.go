package scheduler

import (
	"github.com/gogodjzhu/gogoscrapy/entity"
	"github.com/gogodjzhu/gogoscrapy/utils"
	"github.com/gogodjzhu/gogoscrapy/utils/redisUtil"
)

type DuplicateRemover interface {
	IsDuplicate(request entity.IRequest) bool
	ResetDuplicate()
	GetTotalCount() int
}

type MemDuplicateRemover struct {
	remover *utils.AsyncSet
}

func NewMemDuplicateRemover() *MemDuplicateRemover {
	return &MemDuplicateRemover{remover: utils.NewAsyncSet()}
}

func (this *MemDuplicateRemover) IsDuplicate(request entity.IRequest) bool {
	return !this.remover.Add(request.GetUrl())
}

func (this *MemDuplicateRemover) ResetDuplicate() {
	this.remover.Clear()
}

func (this *MemDuplicateRemover) GetTotalCount() int {
	return this.remover.Size()
}

type RedisDuplicateRemover struct {
	pfkey string
}

func NewRedisDuplicatedRemover(config redisUtil.Config, pfkey string) (*RedisDuplicateRemover, error) {
	if err := redisUtil.Init(config); err != nil {
		return nil, err
	} else {
		return &RedisDuplicateRemover{pfkey: pfkey}, nil
	}
}

func (this *RedisDuplicateRemover) IsDuplicate(request entity.IRequest) bool {
	conn := redisUtil.GetConn()
	defer conn.Close()
	res, err := conn.Do("PFADD", this.pfkey, request.GetUrl())
	if err != nil {
		LOG.Warnf("failed to PFADD to redis so treat this as NotDuplicate, err:%+v", err)
		return false
	}
	return res == 1
}

func (this *RedisDuplicateRemover) ResetDuplicate() {
	conn := redisUtil.GetConn()
	defer conn.Close()
	_, err := conn.Do("DEL", this.pfkey)
	if err != nil {
		LOG.Warnf("failed to DEL HyperLogLog key, err:%+v", err)
	}
}

func (this *RedisDuplicateRemover) GetTotalCount() int {
	conn := redisUtil.GetConn()
	defer conn.Close()
	res, err := conn.Do("PFCOUNT", this.pfkey)
	if err != nil {
		LOG.Warnf("failed to PFCOUNT HyperLogLog key, err:%+v", err)
		return 0
	}
	return res.(int)
}

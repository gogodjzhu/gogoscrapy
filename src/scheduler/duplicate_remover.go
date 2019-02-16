package scheduler

import (
	"gogoscrapy/src/entity"
	"gogoscrapy/src/utils"
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
}

func (this *RedisDuplicateRemover) IsDuplicate(request entity.IRequest) bool {
	panic("implement me")
}

func (this *RedisDuplicateRemover) ResetDuplicate() {
	panic("implement me")
}

func (this *RedisDuplicateRemover) GetTotalCount() int {
	panic("implement me")
}

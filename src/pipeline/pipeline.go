package pipeline

import (
	"gogoscrapy/src/entity"
	"sunteng/commons/log"
)

type IPipeline interface {
	Process(items entity.IResultItems) error
}

type ConsolePipeline struct {
}

func NewConsolePipeline() ConsolePipeline {
	return ConsolePipeline{}
}

func (ConsolePipeline) Process(items entity.IResultItems) error {
	log.Logf("items :%+v", items)
	return nil
}

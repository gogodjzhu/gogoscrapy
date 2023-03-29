package pipeline

import (
	"github.com/gogodjzhu/gogoscrapy/entity"
	"github.com/gogodjzhu/gogoscrapy/utils"
)

var LOG = utils.NewLogger()

type IPipeline interface {
	Process(items entity.IResultItems) error
}

type ConsolePipeline struct {
}

func NewConsolePipeline() ConsolePipeline {
	return ConsolePipeline{}
}

func (ConsolePipeline) Process(items entity.IResultItems) error {
	LOG.Infof("items :%+v", items)
	return nil
}

package pipeline

import (
	"encoding/json"
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
	if bs, err := json.Marshal(items); err == nil {
		LOG.Infof("items:%s", string(bs))
	}
	return nil
}

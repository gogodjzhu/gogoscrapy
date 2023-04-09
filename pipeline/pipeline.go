package pipeline

import (
	"encoding/json"
	"github.com/gogodjzhu/gogoscrapy/entity"
	log "github.com/sirupsen/logrus"
)

type IPipeline interface {
	Pipe(items entity.IResultItems) error
}

type ConsolePipeline struct {
}

func NewConsolePipeline() ConsolePipeline {
	return ConsolePipeline{}
}

func (ConsolePipeline) Pipe(items entity.IResultItems) error {
	if bs, err := json.Marshal(items); err == nil {
		log.Infof("items:%s", string(bs))
	}
	return nil
}

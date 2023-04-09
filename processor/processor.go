package processor

import (
	"github.com/gogodjzhu/gogoscrapy/entity"
)

type IProcessor interface {
	Process(page entity.IPage) error
}

type SimpleProcessor struct {
	urlPattern string
}

func NewSimpleProcessor(urlPattern string) *SimpleProcessor {
	return &SimpleProcessor{urlPattern: urlPattern}
}

func (sp *SimpleProcessor) Process(page entity.IPage) error {
	var links []string
	for _, node := range page.GetHtmlNode().Links().Regex(sp.urlPattern).Nodes() {
		links = append(links, node.Text())
	}
	page.StoreField("url", links)
	page.AddTargetRequestUrls(links...)
	return nil
}

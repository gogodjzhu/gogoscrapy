package processor

import (
	"github.com/gogodjzhu/gogoscrapy/src/entity"
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

func (this *SimpleProcessor) Process(page entity.IPage) error {
	var links []string
	for _, node := range page.GetHtmlNode().Links().Regex(this.urlPattern).Nodes() {
		links = append(links, node.Text())
	}
	page.StoreField("url", links)
	page.AddTargetRequestUrls(links...)
	return nil
}

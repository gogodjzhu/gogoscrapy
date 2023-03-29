package main

import (
	"github.com/gogodjzhu/gogoscrapy"
	"github.com/gogodjzhu/gogoscrapy/downloader"
	"github.com/gogodjzhu/gogoscrapy/entity"
	"github.com/gogodjzhu/gogoscrapy/pipeline"
	"time"
)

type SimpleProcessor struct {
	urlPattern string
}

func NewSimpleProcessor(urlPattern string) *SimpleProcessor {
	return &SimpleProcessor{urlPattern: urlPattern}
}

func (this *SimpleProcessor) Process(page entity.IPage) error {
	var links []string
	//use regex to find links from html.
	for _, node := range page.GetHtmlNode().Links().Regex(this.urlPattern).Nodes() {
		links = append(links, node.Text())
	}
	page.StoreField("url", links) //add links we found to the store, it will be used in pipeline.

	page.AddTargetRequestUrls(links...)
	return nil
}

func main() {
	simplestDemoSpider := gogoscrapy.NewSpider(NewSimpleProcessor("http://.*"))
	simplestDemoSpider.Downloader(downloader.NewSimpleDownloader(10*time.Second, nil))
	simplestDemoSpider.Pipeline(pipeline.NewConsolePipeline())
	simplestDemoSpider.DownloadCoroutineNum(1)
	simplestDemoSpider.DownloadInterval(5 * time.Second)
	simplestDemoSpider.AddStartUrl("http://www.soharp.com")
	simplestDemoSpider.Start()
}

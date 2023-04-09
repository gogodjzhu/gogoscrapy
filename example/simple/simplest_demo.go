package main

import (
	"context"
	"errors"
	"github.com/gogodjzhu/gogoscrapy"
	"github.com/gogodjzhu/gogoscrapy/downloader"
	"github.com/gogodjzhu/gogoscrapy/entity"
	"github.com/gogodjzhu/gogoscrapy/pipeline"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"io"
	"net/url"
	"os"
	"time"
)

type SimpleProcessor struct {
}

func NewSimpleProcessor() *SimpleProcessor {
	return &SimpleProcessor{}
}

func (this *SimpleProcessor) Process(page entity.IPage) error {
	log.Info("processor: ", page.GetRequest().GetUrl())
	originUrl, err := url.Parse(page.GetRequest().GetUrl())
	if err != nil {
		return errors.New("parse url error, url: " + page.GetRequest().GetUrl())
	}
	linkMap := make(map[string]interface{})
	for _, selectable := range page.GetHtmlNode().Links().Nodes() {
		linkStr := selectable.Text()
		u, err := url.Parse(linkStr)
		if err != nil {
			log.Warn("parse url error: ", err)
		}
		if u.Host == "" {
			u.Host = originUrl.Host
			linkStr = originUrl.Scheme + "://" + u.Host + linkStr
		}
		if u.Host != "ssr1.scrape.center" {
			continue
		}
		linkMap[linkStr] = nil
	}
	var links []string
	for link := range linkMap {
		links = append(links, link)
	}
	page.GetResultItems().Put("links", links)
	page.AddTargetRequestUrls(links...)
	return nil
}

func init() {
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat:           "2006-01-02 15:04:05",
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		DisableLevelTruncation:    false,
	})
	rotateWriter, err := rotatelogs.New(
		"logs/logfile.%Y%m%d.log",
		rotatelogs.WithLinkName("logfile"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithRotationSize(1024*1024*1024),
	)
	if err != nil {
		log.Errorf("config local file system for logger error: %v", err)
	}
	writers := []io.Writer{rotateWriter, os.Stdout}
	log.SetOutput(io.MultiWriter(writers...))
}

func main() {
	ctx := context.Background()
	simplestDemoSpider := gogoscrapy.NewSpider(NewSimpleProcessor())
	simplestDemoSpider.Downloader(downloader.NewSimpleDownloader(10*time.Second, nil))
	simplestDemoSpider.Pipeline(pipeline.NewConsolePipeline())
	simplestDemoSpider.DownloadCoroutineNum(1)
	simplestDemoSpider.AddStartUrl("https://ssr1.scrape.center/detail/100")
	simplestDemoSpider.Start(ctx)
}

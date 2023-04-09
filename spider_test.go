package gogoscrapy

import (
	"context"
	"github.com/gogodjzhu/gogoscrapy/downloader"
	"github.com/gogodjzhu/gogoscrapy/entity"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"regexp"
	"testing"
	"time"
)

type testProcessor struct {
}

func (tp *testProcessor) Process(page entity.IPage) error {
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
	var pageLinks []string
	var detailLinks []string
	for link := range linkMap {
		links = append(links, link)
		// regex match page link
		if ok, _ := regexp.MatchString("https://ssr1.scrape.center/page/\\d+", link); ok {
			pageLinks = append(pageLinks, link)
		}
		// regex match detail link
		if ok, _ := regexp.MatchString("https://ssr1.scrape.center/detail/\\w+", link); ok {
			detailLinks = append(detailLinks, link)
		}

	}
	page.GetResultItems().Put("pageLinks", pageLinks)
	page.GetResultItems().Put("detailLinks", detailLinks)
	page.AddTargetRequestUrls(links...)
	return nil
}

type testPipeline struct {
	pageLinkCnt   map[string]int
	detailLinkCnt map[string]int
}

func (tp *testPipeline) Pipe(items entity.IResultItems) error {
	if items.Get("detailLinks") != nil {
		for _, link := range items.Get("detailLinks").([]string) {
			_, ok := tp.detailLinkCnt[link]
			if !ok {
				tp.detailLinkCnt[link] = 0
			}
			tp.detailLinkCnt[link]++
			log.Infof("detailLink: %s, cnt: %d", link, tp.detailLinkCnt[link])
		}
	}
	if items.Get("pageLinks") != nil {
		for _, link := range items.Get("pageLinks").([]string) {
			_, ok := tp.pageLinkCnt[link]
			if !ok {
				tp.pageLinkCnt[link] = 0
			}
			tp.pageLinkCnt[link]++
			log.Infof("pageLink: %s, cnt: %d", link, tp.pageLinkCnt[link])
		}
	}
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

	go func() {
		log.Fatalln(http.ListenAndServe(":9876", nil))
	}()
}

func TestSpider_Basic(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	tp := &testPipeline{
		detailLinkCnt: make(map[string]int),
		pageLinkCnt:   make(map[string]int),
	}
	go func() {
		for {
			select {
			case <-time.After(60 * time.Second):
				cancel()
				return
			default:
				if len(tp.detailLinkCnt) >= 100 {
					cancel()
					return
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	simplestDemoSpider := NewSpider(&testProcessor{})
	simplestDemoSpider.Downloader(downloader.NewSimpleDownloader(10*time.Second, nil))
	simplestDemoSpider.DownloadCoroutineNum(10)
	simplestDemoSpider.Pipeline(tp)
	simplestDemoSpider.AddStartUrl("https://ssr1.scrape.center/detail/100")
	simplestDemoSpider.Start(ctx) // start the spider, this is a blocking call until ctx canceled

	if len(tp.detailLinkCnt) != 100 {
		t.Errorf("test failed, detailLinkCnt: %d", len(tp.detailLinkCnt))
	}
}

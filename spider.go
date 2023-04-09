package gogoscrapy

import (
	"context"
	"errors"
	"github.com/gogodjzhu/gogoscrapy/downloader"
	"github.com/gogodjzhu/gogoscrapy/entity"
	buzzErr "github.com/gogodjzhu/gogoscrapy/err"
	"github.com/gogodjzhu/gogoscrapy/pipeline"
	"github.com/gogodjzhu/gogoscrapy/processor"
	"github.com/gogodjzhu/gogoscrapy/scheduler"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type ISpider interface {
	IApp
	Downloader(downloader downloader.IDownloader) ISpider
	Pipeline(pipeline pipeline.IPipeline) ISpider
	Processor(processor processor.IProcessor) ISpider
	Scheduler(scheduler scheduler.IScheduler) ISpider
	DownloadCoroutineNum(num int) ISpider
	AddStartUrl(urls ...string) ISpider
	AddStartRequest(request entity.IRequest) ISpider
	MaxDownloadRetryTime(rt int) ISpider
}

type spider struct {
	downloader downloader.IDownloader
	pipelines  []pipeline.IPipeline
	processor  processor.IProcessor
	scheduler  scheduler.IScheduler

	startRequests []entity.IRequest

	downloaderCoroutineNum int
	processorCoroutineNum  int
	maxDownloadRetryTime   int // download retry time, -1 means retry forever

	downloadingNum int32             // downloading num
	pageChan       chan entity.IPage // page channel between downloader and processor

	wg *sync.WaitGroup // wait group for all goroutines
}

func NewSpider(proc processor.IProcessor) ISpider {
	return &spider{
		processor: proc,

		pageChan: make(chan entity.IPage, 100),

		wg: &sync.WaitGroup{},
	}
}

func (sp *spider) Downloader(downloader downloader.IDownloader) ISpider {
	sp.downloader = downloader
	return sp
}

func (sp *spider) Pipeline(pipeline pipeline.IPipeline) ISpider {
	sp.pipelines = append(sp.pipelines, pipeline)
	return sp
}

func (sp *spider) Processor(processor processor.IProcessor) ISpider {
	sp.processor = processor
	return sp
}

func (sp *spider) Scheduler(scheduler scheduler.IScheduler) ISpider {
	sp.scheduler = scheduler
	return sp
}

func (sp *spider) DownloadCoroutineNum(num int) ISpider {
	sp.downloaderCoroutineNum = num
	return sp
}

func (sp *spider) AddStartUrl(urls ...string) ISpider {
	if sp.startRequests == nil {
		sp.startRequests = make([]entity.IRequest, 0)
	}
	for _, url := range urls {
		sp.startRequests = append(sp.startRequests, &entity.Request{Url: url, Method: http.MethodGet})
	}
	return sp
}

func (sp *spider) AddStartRequest(request entity.IRequest) ISpider {
	if sp.startRequests == nil {
		sp.startRequests = make([]entity.IRequest, 0)
	}
	sp.startRequests = append(sp.startRequests, request)
	return sp
}

func (sp *spider) MaxDownloadRetryTime(rt int) ISpider {
	sp.maxDownloadRetryTime = rt
	return sp
}

func (sp *spider) Start(ctx context.Context) {
	sp.init(ctx)
	if err := sp.doScrapy(ctx); err != nil {
		log.Errorf("failed scrapy, err:%+v", err)
	}
}

// actually scrapy page
func (sp *spider) doScrapy(ctx context.Context) error {

	//parallel downloader
	for i := 0; i < sp.downloaderCoroutineNum; i++ {
		sp.wg.Add(1)
		go func(index int, ctx context.Context, wg *sync.WaitGroup) {
			for {
				select {
				case <-ctx.Done():
					log.Infof("downloader[%d] is canceled.", index)
					wg.Done()
					return
				case req := <-sp.scheduler.PullChan():
					page := sp.download(req)
					if page != nil {
						sp.pageChan <- page
					}
				}
			}
		}(i, ctx, sp.wg)
	}

	//parallel processor
	for i := 0; i < sp.processorCoroutineNum; i++ {
		sp.wg.Add(1)
		go func(index int, ctx context.Context, wg *sync.WaitGroup) {
			for {
				select {
				case <-ctx.Done():
					log.Infof("processor[%d] is canceled.", index)
					wg.Done()
					return
				case page := <-sp.pageChan:
					if err := sp.processor.Process(page); err != nil {
						if errors.Is(err, buzzErr.RetryAbleError) {
							sp.reSchedule(page.GetRequest())
						} else {
							log.Errorf("processor failed to process, err:%+v", err)
						}
						continue
					}
					if !page.GetResultItems().IsSkip() {
						for _, pipe := range sp.pipelines {
							if err := pipe.Pipe(page.GetResultItems()); err != nil {
								log.Errorf("pipeline[%+v] failed to process, err:%+v", pipe, err)
								continue
							}
						}
					}
					for _, req := range page.GetTargetRequests() {
						sp.scheduler.PushChan() <- req
					}
				}
			}
		}(i, ctx, sp.wg)
	}

	sp.wg.Wait()
	return nil
}

func (sp *spider) reSchedule(req entity.IRequest) {
	var nextTimes int
	times := req.GetExtras()[entity.CycleTriedTimes]
	if times == nil {
		nextTimes = 1
	} else {
		nextTimes = times.(int) + 1
	}
	req.PutExtra(entity.CycleTriedTimes, nextTimes)
	if sp.maxDownloadRetryTime == -1 || nextTimes <= sp.maxDownloadRetryTime {
		sp.scheduler.PushChan() <- req
	} else {
		log.Warn("give up download for retry time over")
	}
}

//init spider and set the params
func (sp *spider) init(ctx context.Context) {
	if sp.downloader == nil {
		sp.downloader = downloader.NewSimpleDownloader(10*time.Second, nil)
	}
	if sp.downloaderCoroutineNum < 1 {
		sp.downloaderCoroutineNum = 1
	}
	if sp.processorCoroutineNum < 1 {
		sp.processorCoroutineNum = 1
	}
	if sp.pipelines == nil {
		sp.pipelines = append(sp.pipelines, pipeline.NewConsolePipeline())
	}
	if sp.scheduler == nil {
		sp.scheduler = scheduler.NewMemScheduler(ctx, sp.wg)
	}
	if sp.maxDownloadRetryTime < 1 {
		sp.maxDownloadRetryTime = 3
	}
	for _, request := range sp.startRequests {
		sp.scheduler.PushChan() <- request
	}
	sp.startRequests = nil //clear object
}

func (sp *spider) download(req entity.IRequest) entity.IPage {
	atomic.AddInt32(&sp.downloadingNum, 1)
	defer atomic.AddInt32(&sp.downloadingNum, -1)
	page, err := sp.downloader.Download(req)
	if err != nil {
		sp.onDownloadFail(req, err)
		sp.reSchedule(req)
		return nil
	}
	sp.onDownloadSuccess(req, page)
	return page
}

func (sp *spider) onDownloadSuccess(req entity.IRequest, page entity.IPage) {
	codes := sp.downloader.GetAcceptStatus()
	for _, code := range codes {
		if code == page.GetStatusCode() {
			proxy, ok := req.GetExtras()[entity.AssignedProxy]
			if ok {
				sp.downloader.GetProxyFactory().ReturnProxy(proxy.(downloader.IProxy))
				delete(req.GetExtras(), entity.AssignedProxy)
			}
			return
		}
	}
	log.Debugf("download failed, status code:%d, url:%s", page.GetStatusCode(), req.GetUrl())
	sp.reSchedule(req)
}

func (sp *spider) onDownloadFail(req entity.IRequest, err error) {
	log.Debugf("download failed, err:%+v, url:%s", err, req.GetUrl())
}

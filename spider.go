package gogoscrapy

import (
	"github.com/gogodjzhu/gogoscrapy/downloader"
	"github.com/gogodjzhu/gogoscrapy/entity"
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
}

const (
	StatInit = iota
	StatRunning
	StatStopped
)

type spider struct {
	downloader downloader.IDownloader

	pipelines []pipeline.IPipeline
	processor processor.IProcessor
	scheduler scheduler.IScheduler

	stat          atomic.Value
	startRequests []entity.IRequest

	downloaderCoroutineNum int
	processorChan          chan entity.IPage
	processorCoroutineNum  int
	retryTime              int
	downloadInterval       time.Duration
}

func NewSpider(proc processor.IProcessor) *spider {
	return &spider{
		processor: proc,
	}
}

func (sp *spider) Downloader(downloader downloader.IDownloader) *spider {
	sp.downloader = downloader
	return sp
}

func (sp *spider) Pipeline(pipeline pipeline.IPipeline) *spider {
	sp.pipelines = append(sp.pipelines, pipeline)
	return sp
}

func (sp *spider) Processor(processor processor.IProcessor) *spider {
	sp.processor = processor
	return sp
}

func (sp *spider) Scheduler(scheduler scheduler.IScheduler) *spider {
	sp.scheduler = scheduler
	return sp
}

func (sp *spider) DownloadCoroutineNum(num int) *spider {
	sp.downloaderCoroutineNum = num
	return sp
}

func (sp *spider) AddStartUrl(urls ...string) {
	if sp.startRequests == nil {
		sp.startRequests = make([]entity.IRequest, 0)
	}
	for _, url := range urls {
		sp.startRequests = append(sp.startRequests, &entity.Request{Url: url, Method: http.MethodGet})
	}
}

func (sp *spider) AddStartRequest(request entity.IRequest) {
	if sp.startRequests == nil {
		sp.startRequests = make([]entity.IRequest, 0)
	}
	sp.startRequests = append(sp.startRequests, request)
}

func (sp *spider) RetryTime(rt int) *spider {
	sp.retryTime = rt
	return sp
}

func (sp *spider) DownloadInterval(di time.Duration) *spider {
	sp.downloadInterval = di
	return sp
}

func (sp *spider) Start() {
	if sp.stat.Load() == StatRunning {
		panic("spider is already running.")
	}
	sp.init()
	sp.stat.Store(StatRunning)
	if sp.stat.Load() == StatRunning {
		if err := sp.doScrapy(); err != nil {
			log.Errorf("failed scrapy, err:%+v", err)
		}
	}
}

// actually scrapy page
func (sp *spider) doScrapy() error {
	wg := sync.WaitGroup{}
	var downloadingNum int32

	//exit monitor
	go func() {
		for {
			time.Sleep(30 * time.Second)
			log.Infof("task remain in scheduler: %d, downloading: %d", sp.scheduler.Size(), downloadingNum)
			//if downloading task is none then shutdown
			if sp.scheduler.IsClose() || (sp.scheduler.Size() < 1 && downloadingNum < 1) {
				//wait and double check
				time.Sleep(30 * time.Second)
				if sp.scheduler.IsClose() || (sp.scheduler.Size() < 1 && downloadingNum < 1) {
					log.Infof("no more task to download, shutdown it gracefully.")
					sp.Shutdown()
				}
			}
		}
	}()

	//parallel downloader
	for i := 0; i < sp.downloaderCoroutineNum; i++ {
		wg.Add(1)
		go func(index int) {
			for {
				req := sp.scheduler.Poll()
				if req == nil {
					if sp.scheduler.IsClose() {
						break
					}
					//wait next request
					time.Sleep(3 * time.Second)
					log.Debugf("wait for next request")
					continue
				}
				atomic.AddInt32(&downloadingNum, 1)
				if page, err := sp.downloader.Download(req); err != nil {
					sp.onDownloadFail(req, err)
					sp.doRetry(req)
				} else {
					sp.onDownloadSuccess(req, page)
					sp.processorChan <- page
				}
				atomic.AddInt32(&downloadingNum, -1)
				time.Sleep(sp.downloadInterval)
			}
			log.Infof("downloader[%d] end.", index)
			wg.Done()
		}(i)
	}

	//parallel processor
	for i := 0; i < sp.processorCoroutineNum; i++ {
		wg.Add(1)
		go func(index int) {
			for page := range sp.processorChan {
				if err := sp.processor.Process(page); err != nil {
					log.Errorf("processor failed to process, err:%+v", err)
					sp.doRetry(page.GetRequest())
					continue
				}
				for _, req := range page.GetTargetRequests() {
					sp.scheduler.Push(req)
				}
				if !page.GetResultItems().IsSkip() {
					for _, pipe := range sp.pipelines {
						if err := pipe.Process(page.GetResultItems()); err != nil {
							log.Errorf("pipeline[%+v] failed to process, err:%+v", pipe, err)
							continue
						}
					}
				}
			}
			log.Infof("processor[%d] end.", index)
			wg.Done()
		}(i)
	}

	wg.Wait()
	return nil
}

func (sp *spider) doRetry(req entity.IRequest) {
	var nextTimes int
	times := req.GetExtras()[entity.CycleTriedTimes]
	if times == nil {
		nextTimes = 1
	} else {
		nextTimes = times.(int) + 1
	}
	req.PutExtra(entity.CycleTriedTimes, nextTimes)
	if nextTimes <= sp.retryTime {
		sp.scheduler.Push(req)
	} else {
		log.Warn("give up download for retry time over")
	}
}

func (sp *spider) Shutdown() {
	sp.stat.Store(StatInit)
	log.Infof("shutdown...")
	tryClose(sp.downloader)
	tryClose(sp.scheduler)
	close(sp.processorChan)
	tryClose(sp.processor)
	for pipe := range sp.pipelines {
		tryClose(pipe)
	}
	sp.stat.Store(StatStopped)
}

func tryClose(closeable interface{}) {
	if c, ok := closeable.(entity.Closeable); ok {
		if err := c.Close(); err != nil {
			log.Warnf("failed to close %+v, err:%+v", c, err)
		}
	}
}

func (sp *spider) IsShutdown() bool {
	return sp.stat.Load().(int) == StatStopped
}

//init spider and set the params
func (sp *spider) init() {
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
		sp.scheduler = scheduler.NewQueueScheduler()
	}
	if sp.retryTime < 1 {
		sp.retryTime = 3
	}
	if sp.downloadInterval < 1 {
		sp.downloadInterval = 10 * time.Second
	}
	for _, request := range sp.startRequests {
		sp.scheduler.Push(request)
	}
	sp.startRequests = nil //clear object
	sp.processorChan = make(chan entity.IPage, 1)
	sp.stat.Store(StatInit)
}

func (sp *spider) onDownloadSuccess(req entity.IRequest, page entity.IPage) {
	codes := sp.downloader.GetAcceptStatus()
	for _, code := range codes {
		if code == page.GetStatusCode() {
			proxy := req.GetExtras()[entity.AssignedProxy]
			sp.downloader.GetProxyFactory().ReturnProxy(proxy.(downloader.IProxy))
			delete(req.GetExtras(), entity.AssignedProxy)
			return
		}
	}
	log.Debugf("download failed, status code:%d, url:%s", page.GetStatusCode(), req.GetUrl())
	sp.doRetry(req)
}

func (sp *spider) onDownloadFail(req entity.IRequest, err error) {
	log.Debugf("download failed, err:%+v, url:%s", err, req.GetUrl())
}

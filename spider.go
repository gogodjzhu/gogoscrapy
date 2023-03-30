package gogoscrapy

import (
	"github.com/gogodjzhu/gogoscrapy/downloader"
	"github.com/gogodjzhu/gogoscrapy/entity"
	"github.com/gogodjzhu/gogoscrapy/pipeline"
	"github.com/gogodjzhu/gogoscrapy/processor"
	"github.com/gogodjzhu/gogoscrapy/scheduler"
	"github.com/gogodjzhu/gogoscrapy/utils"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var LOG = utils.NewLogger()

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

func (this *spider) Downloader(downloader downloader.IDownloader) *spider {
	this.downloader = downloader
	return this
}

func (this *spider) Pipeline(pipeline pipeline.IPipeline) *spider {
	this.pipelines = append(this.pipelines, pipeline)
	return this
}

func (this *spider) Processor(processor processor.IProcessor) *spider {
	this.processor = processor
	return this
}

func (this *spider) Scheduler(scheduler scheduler.IScheduler) *spider {
	this.scheduler = scheduler
	return this
}

func (this *spider) DownloadCoroutineNum(num int) *spider {
	this.downloaderCoroutineNum = num
	return this
}

func (this *spider) AddStartUrl(urls ...string) {
	if this.startRequests == nil {
		this.startRequests = make([]entity.IRequest, 0)
	}
	for _, url := range urls {
		this.startRequests = append(this.startRequests, &entity.Request{Url: url, Method: http.MethodGet})
	}
}

func (this *spider) AddStartRequest(request entity.IRequest) {
	if this.startRequests == nil {
		this.startRequests = make([]entity.IRequest, 0)
	}
	this.startRequests = append(this.startRequests, request)
}

func (this *spider) RetryTime(rt int) *spider {
	this.retryTime = rt
	return this
}

func (this *spider) DownloadInterval(di time.Duration) *spider {
	this.downloadInterval = di
	return this
}

func (this *spider) Start() {
	if this.stat.Load() == StatRunning {
		panic("spider is already running.")
	}
	this.init()
	this.stat.Store(StatRunning)
	if this.stat.Load() == StatRunning {
		if err := this.doScrapy(); err != nil {
			LOG.Errorf("failed scrapy, err:%+v", err)
		}
	}
}

// actually scrapy page
func (this *spider) doScrapy() error {
	wg := sync.WaitGroup{}
	var downloadingNum int32

	//exit monitor
	go func() {
		for {
			time.Sleep(30 * time.Second)
			LOG.Infof("task remain in scheduler: %d, downloading: %d", this.scheduler.Size(), downloadingNum)
			//if downloading task is none then shutdown
			if this.scheduler.IsClose() || (this.scheduler.Size() < 1 && downloadingNum < 1) {
				//wait and double check
				time.Sleep(30 * time.Second)
				if this.scheduler.IsClose() || (this.scheduler.Size() < 1 && downloadingNum < 1) {
					LOG.Infof("no more task to download, shutdown it gracefully.")
					this.Shutdown()
				}
			}
		}
	}()

	//parallel downloader
	for i := 0; i < this.downloaderCoroutineNum; i++ {
		wg.Add(1)
		go func(index int) {
			for {
				req := this.scheduler.Poll()
				if req == nil {
					if this.scheduler.IsClose() {
						break
					}
					//wait next request
					time.Sleep(3 * time.Second)
					LOG.Debugf("wait for next request")
					continue
				}
				atomic.AddInt32(&downloadingNum, 1)
				if page, err := this.downloader.Download(req); err != nil {
					LOG.Warnf("failed to download, err：%+v", err)
					this.doRetry(req)
				} else {
					LOG.Debugf("success download, url:%+s", req.GetUrl())
					this.processorChan <- page
				}
				atomic.AddInt32(&downloadingNum, -1)
				time.Sleep(this.downloadInterval)
			}
			LOG.Infof("downloader[%d] end.", index)
			wg.Done()
		}(i)
	}

	//parallel processor
	for i := 0; i < this.processorCoroutineNum; i++ {
		wg.Add(1)
		go func(index int) {
			for page := range this.processorChan {
				if err := this.processor.Process(page); err != nil {
					LOG.Errorf("processor failed to process, err:%+v", err)
					this.doRetry(page.GetRequest())
					continue
				}
				for _, req := range page.GetTargetRequests() {
					this.scheduler.Push(req)
				}
				if !page.GetResultItems().IsSkip() {
					for _, pipe := range this.pipelines {
						if err := pipe.Process(page.GetResultItems()); err != nil {
							LOG.Errorf("pipeline[%+v] failed to process, err:%+v", pipe, err)
							continue
						}
					}
				}
			}
			LOG.Infof("processor[%d] end.", index)
			wg.Done()
		}(i)
	}

	wg.Wait()
	return nil
}

func (this *spider) doRetry(req entity.IRequest) {
	var nextTimes int
	times := req.GetExtras()[entity.CycleTriedTimes]
	if times == nil {
		nextTimes = 1
	} else {
		nextTimes = times.(int) + 1
	}
	req.PutExtra(entity.CycleTriedTimes, nextTimes)
	if nextTimes <= this.retryTime {
		this.scheduler.Push(req)
	} else {
		LOG.Warn("give up download for retry time over")
	}
}

func (this *spider) Shutdown() {
	this.stat.Store(StatInit)
	LOG.Infof("shutdown...")
	tryClose(this.downloader)
	tryClose(this.scheduler)
	close(this.processorChan)
	tryClose(this.processor)
	for pipe := range this.pipelines {
		tryClose(pipe)
	}
	this.stat.Store(StatStopped)
}

func tryClose(closeable interface{}) {
	if c, ok := closeable.(entity.Closeable); ok {
		if err := c.Close(); err != nil {
			LOG.Warnf("failed to close %+v, err:%+v", c, err)
		}
	}
}

func (this *spider) IsShutdown() bool {
	return this.stat.Load().(int) == StatStopped
}

//init spider and set the params
func (this *spider) init() {
	if this.downloader == nil {
		this.downloader = downloader.NewSimpleDownloader(10*time.Second, nil)
	}
	if this.downloaderCoroutineNum < 1 {
		this.downloaderCoroutineNum = 1
	}
	if this.processorCoroutineNum < 1 {
		this.processorCoroutineNum = 1
	}
	if this.pipelines == nil {
		this.pipelines = append(this.pipelines, pipeline.NewConsolePipeline())
	}
	if this.scheduler == nil {
		this.scheduler = scheduler.NewQueueScheduler()
	}
	if this.retryTime < 1 {
		this.retryTime = 3
	}
	if this.downloadInterval < 1 {
		this.downloadInterval = 10 * time.Second
	}
	for _, request := range this.startRequests {
		this.scheduler.Push(request)
	}
	this.startRequests = nil //clear object
	this.processorChan = make(chan entity.IPage, 1)
	this.stat.Store(StatInit)
}

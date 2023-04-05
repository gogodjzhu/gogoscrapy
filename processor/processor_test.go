package processor

import (
	"github.com/gogodjzhu/gogoscrapy/downloader"
	"github.com/gogodjzhu/gogoscrapy/entity"
	"testing"
	"time"
)

func TestSimpleProcessor_Process(t *testing.T) {
	processor := NewSimpleProcessor("http://.*")
	simpleDownloader := downloader.NewSimpleDownloader(10*time.Second, nil)

	req := entity.NewGetRequest("http://www.baidu.com")
	page, err := simpleDownloader.Download(req)
	if err != nil {
		simpleDownloader.OnError(req, err)
	} else {
		simpleDownloader.OnSuccess(req)
		if err := processor.Process(page); err != nil {
			t.Errorf("TestSimpleProcessor_Process failed, err:%+v", err)
		}
		if len(page.GetTargetRequests()) < 1 {
			t.Errorf("TestSimpleProcessor_Process may failed for the targetRequest is empty, " +
				"please check the processor")
		}
	}
}

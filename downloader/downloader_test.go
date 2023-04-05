package downloader

import (
	"github.com/gogodjzhu/gogoscrapy/entity"
	"testing"
	"time"
)

func TestSimpleDownloader_Download(t *testing.T) {
	simpleDownloader := NewSimpleDownloader(10*time.Second, nil)
	req, err := entity.NewGetRequest("http://gogodjzhu.com")
	if err != nil {
		t.FailNow()
	}
	_, err = simpleDownloader.Download(req)
	if err != nil {
		t.FailNow()
	}
}

func TestSimpleDownload_IllegalRequest(t *testing.T) {
	_, err := entity.NewGetRequest("http:/gogodjzhu.com")
	if err == nil {
		t.FailNow()
	}
}

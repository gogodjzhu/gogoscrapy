package downloader

import (
	"github.com/gogodjzhu/gogoscrapy/src/entity"
	"testing"
	"time"
)

func TestSimpleDownloader_Download(t *testing.T) {
	simpleDownloader := NewSimpleDownloader(10*time.Second, nil)
	req := entity.NewRequest("http://gogodjzhu.com")
	_, err := simpleDownloader.Download(req)
	if err != nil {
		simpleDownloader.OnError(req, err)
		t.Errorf("TestSimpleDownloader_Download failed, err:%+v", err)
	} else {
		simpleDownloader.OnSuccess(req)
	}
}

func TestSimpleDownloader_Download_With_Proxy(t *testing.T) {
	proxyFactory, err := NewFileProxyFactory("../../files/proxy_list")
	if err != nil {
		panic(err)
	}
	simpleDownloader := NewSimpleDownloader(10*time.Second, proxyFactory)
	req := entity.NewRequest("http://www.gogodjzhu.com").SetUseProxy(true)
	_, err = simpleDownloader.Download(req)
	if err != nil {
		simpleDownloader.OnError(req, err)
		t.Errorf("TestSimpleDownloader_Download_With_Proxy failed, err:%+v", err)
	} else {
		simpleDownloader.OnSuccess(req)
	}
}

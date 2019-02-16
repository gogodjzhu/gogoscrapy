package downloader

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"gogoscrapy/src/entity"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sunteng/commons/log"
	"time"
)

type IDownloader interface {
	Download(request entity.IRequest) (entity.IPage, error)
	SetDownloadTimeout(dt time.Duration)
	OnSuccess(request entity.IRequest)
	OnError(request entity.IRequest, err error)
}

type simpleDownloader struct {
	downloadTimeout time.Duration
	proxyFactory    IProxyFactory
}

func NewSimpleDownloader(downloadTimeout time.Duration, provider IProxyFactory) *simpleDownloader {
	return &simpleDownloader{downloadTimeout: downloadTimeout, proxyFactory: provider}
}

func (this *simpleDownloader) Download(request entity.IRequest) (entity.IPage, error) {
	client, proxy, err := this.getHttpRequest(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(request.GetMethod(), request.GetUrl(), strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	headers := request.GetHeaders()
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, strings.Join(v, ";"))
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	if proxy != nil {
		this.proxyFactory.ReturnProxy(proxy)
	}
	return entity.NewPage(request, doc, request.GetCharset(), resp.StatusCode, resp.Header, false), nil
}

func (this *simpleDownloader) SetDownloadTimeout(dt time.Duration) {
	this.downloadTimeout = dt
}

func (this *simpleDownloader) OnSuccess(request entity.IRequest) {
	log.Debugf("success download page, url:%s", request.GetUrl())
}

func (this *simpleDownloader) OnError(request entity.IRequest, err error) {
	log.Warnf("failed to download page, request:%+v, err:%+v", request, err)
}

func (this *simpleDownloader) getHttpRequest(request entity.IRequest) (*http.Client, IProxy, error) {
	client := http.Client{Timeout: this.downloadTimeout}
	var proxy IProxy
	var err error
	var urlProxy *url.URL
	if request.IsUseProxy() {
		if this.proxyFactory == nil {
			return nil, nil, errors.New("request want to use proxy but proxyProvider is nil")
		}
		proxy, err = this.proxyFactory.GetProxy()
		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("failed to get proxy, error:%+v", err))
		}
		urlProxy, err = url.Parse("http://" + proxy.GetHost() + ":" + strconv.Itoa(proxy.GetPort()))
		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("failed to get proxy, error:%+v", err))
		}
	}

	client.Transport = &http.Transport{
		Proxy: http.ProxyURL(urlProxy),
	}
	return &client, proxy, nil
}

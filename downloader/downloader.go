package downloader

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	ent "github.com/gogodjzhu/gogoscrapy/entity"
	"github.com/gogodjzhu/gogoscrapy/utils"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var LOG = utils.NewLogger()

type IDownloader interface {
	Download(request ent.IRequest) (ent.IPage, error)
	SetDownloadTimeout(dt time.Duration)
	OnSuccess(request ent.IRequest)
	OnError(request ent.IRequest, err error)
	Validate(page ent.IPage) (bool, string)
}

type SimpleDownloader struct {
	downloadTimeout time.Duration
	proxyFactory    IProxyFactory
}

func NewSimpleDownloader(downloadTimeout time.Duration, provider IProxyFactory) *SimpleDownloader {
	return &SimpleDownloader{downloadTimeout: downloadTimeout, proxyFactory: provider}
}

func (rd *SimpleDownloader) Download(request ent.IRequest) (ent.IPage, error) {
	client, proxy, err := rd.getHttpRequest(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(request.GetMethod(), request.GetUrl(), strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36")
	req.Header.Set("Accept", "*/*")
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
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	rawText, err := getRawText(resp.Header, doc)
	if err != nil {
		LOG.Warnf("failed to get html from document, err:%+v", err)
		return nil, err
	}
	defer resp.Body.Close()
	page := ent.NewPage(request, doc, getCharset(resp.Header), resp.StatusCode, resp.Header, rawText, false)
	if valid, msg := rd.Validate(page); !valid {
		return nil, errors.New(msg)
	}
	if proxy != nil {
		rd.proxyFactory.ReturnProxy(proxy)
	}
	return page, nil
}

func getRawText(header http.Header, doc *goquery.Document) (string, error) {
	if len(strings.TrimSpace(header.Get("Content-Type"))) > 0 {
		//eg. application/json; charset=utf-8
		contentTypeStr := header.Get("Content-Type")
		ctArr := strings.Split(contentTypeStr, ";")
		for _, pair := range ctArr {
			switch {
			case strings.TrimSpace(pair) == "application/json":
				return doc.Text(), nil
			case strings.TrimSpace(pair) == "text/xml":
				return doc.Find("body").Html()
			case strings.TrimSpace(pair) == "text/html":
				return doc.Html()
			}
		}
	}
	return doc.Html()
}

func getCharset(header http.Header) string {
	if len(strings.TrimSpace(header.Get("Content-Type"))) > 0 {
		//eg. application/json; charset=utf-8
		contentTypeStr := header.Get("Content-Type")
		ctArr := strings.Split(contentTypeStr, ";")
		for _, pair := range ctArr {
			kvArr := strings.Split(strings.TrimSpace(pair), "=")
			if len(kvArr) == 2 && kvArr[0] == "charset" {
				return kvArr[1]
			}
		}
	}
	return ""
}

func (rd *SimpleDownloader) SetDownloadTimeout(dt time.Duration) {
	rd.downloadTimeout = dt
}

func (rd *SimpleDownloader) OnSuccess(request ent.IRequest) {
	LOG.Debugf("success download page, url:%s", request.GetUrl())
}

func (rd *SimpleDownloader) OnError(request ent.IRequest, err error) {
	LOG.Warnf("failed to download page, request:%+v, err:%+v", request, err)
}

func (rd *SimpleDownloader) Validate(ent.IPage) (bool, string) {
	return true, "ok"
}

func (rd *SimpleDownloader) getHttpRequest(request ent.IRequest) (*http.Client, IProxy, error) {
	client := http.Client{Timeout: rd.downloadTimeout}
	var proxy IProxy
	var err error
	var urlProxy *url.URL
	if request.IsUseProxy() {
		if rd.proxyFactory == nil {
			return nil, nil, errors.New("request want to use proxy but proxyProvider is nil")
		}
		proxy, err = rd.proxyFactory.GetProxy()
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

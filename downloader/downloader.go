package downloader

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	ent "github.com/gogodjzhu/gogoscrapy/entity"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type IDownloader interface {
	Download(request ent.IRequest) (ent.IPage, error)
	SetDownloadTimeout(dt time.Duration)
	GetDownloadTimeout() time.Duration
	SetProxyFactory(provider IProxyFactory)
	GetProxyFactory() IProxyFactory
	SetAcceptStatus(status []int)
	GetAcceptStatus() []int
}

type SimpleDownloader struct {
	downloadTimeout time.Duration
	proxyFactory    IProxyFactory
	acceptStatus    []int
}

func NewSimpleDownloader(downloadTimeout time.Duration, provider IProxyFactory) *SimpleDownloader {
	return &SimpleDownloader{
		downloadTimeout: downloadTimeout,
		proxyFactory:    provider,
		acceptStatus:    []int{200},
	}
}

func (rd *SimpleDownloader) Download(request ent.IRequest) (ent.IPage, error) {
	var client *http.Client
	var req *http.Request
	var err error
	if client, err = rd.wrapClient(request); err != nil {
		return nil, err
	}
	if request.GetBody() != nil {
		if req, err = http.NewRequest(request.GetMethod(), request.GetUrl(), bytes.NewReader(request.GetBody())); err != nil {
			return nil, err
		}
	} else {
		if req, err = http.NewRequest(request.GetMethod(), request.GetUrl(), nil); err != nil {
			return nil, err
		}
	}

	// set headers
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
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	rawText, err := getRawText(resp.Header, doc)
	if err != nil {
		log.Warnf("failed to get html from document, err:%+v", err)
		return nil, err
	}
	return ent.NewPage(request, doc, getCharset(resp.Header), resp.StatusCode, resp.Header, rawText, false), nil
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

func (rd *SimpleDownloader) GetDownloadTimeout() time.Duration {
	return rd.downloadTimeout
}

func (rd *SimpleDownloader) SetProxyFactory(provider IProxyFactory) {
	rd.proxyFactory = provider
}

func (rd *SimpleDownloader) GetProxyFactory() IProxyFactory {
	return rd.proxyFactory
}

func (rd *SimpleDownloader) SetAcceptStatus(status []int) {
	rd.acceptStatus = status
}

func (rd *SimpleDownloader) GetAcceptStatus() []int {
	return rd.acceptStatus
}

func (rd *SimpleDownloader) wrapClient(request ent.IRequest) (*http.Client, error) {
	client := http.Client{Timeout: rd.downloadTimeout}
	if !request.IsUseProxy() {
		return &client, nil
	}
	var proxy IProxy
	var err error
	if rd.proxyFactory == nil {
		log.Warn("request want to use proxy but proxyProvider is nil, use default client")
		return &client, nil
	}
	if proxy, err = rd.proxyFactory.GetProxy(); err != nil {
		log.Warn("failed to get proxy, use default client")
		return &client, nil
	}
	request.PutExtra(ent.AssignedProxy, proxy)
	client.Transport = proxy.GetTransport()
	return &client, nil
}

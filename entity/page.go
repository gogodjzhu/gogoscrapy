package entity

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gogodjzhu/gogoscrapy/selector"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

//represent a html page, all methods are not thread safe.
type IPage interface {
	GetRequest() IRequest
	GetDocument() *goquery.Document
	GetHtmlNode() *selector.HtmlNode
	GetCharset() string
	GetStatusCode() int
	AddTargetRequestUrls(url ...string)
	AddTargetRequests(request ...IRequest)
	GetTargetRequests() []IRequest
	GetRespHeaders() map[string][]string
	GetUrl() *selector.PlainText
	StoreField(key string, obj interface{})
	GetResultItems() IResultItems
	GetRawText() string
}

type Page struct {
	Request        IRequest
	Document       *goquery.Document
	HtmlRootNode   *selector.HtmlNode
	Charset        string
	Status         int
	TargetRequests []IRequest
	RespHeaders    map[string][]string
	Url            *selector.PlainText
	rawText        string
	resultItems    IResultItems
	domain         string
}

func NewPage(request IRequest, document *goquery.Document, charset string, status int,
	respHeader map[string][]string, rawText string, skip bool) *Page {
	return &Page{
		Request:        request,
		Document:       document,
		HtmlRootNode:   &selector.HtmlNode{Elements: document.Nodes},
		Charset:        charset,
		Status:         status,
		TargetRequests: make([]IRequest, 0),
		RespHeaders:    respHeader,
		Url:            &selector.PlainText{SourceTexts: []string{request.GetUrl()}},
		resultItems:    NewResultItems(skip),
		domain: func() string {
			urlStr := request.GetUrl()
			u, err := url.Parse(urlStr)
			if err != nil {
				panic(err)
			}
			return u.Scheme + "://" + u.Host
		}(),
		rawText: rawText,
	}
}

func (pg *Page) GetRequest() IRequest {
	return pg.Request
}

func (pg *Page) GetDocument() *goquery.Document {
	return pg.Document
}

func (pg *Page) GetHtmlNode() *selector.HtmlNode {
	return pg.HtmlRootNode
}

func (pg *Page) GetCharset() string {
	return pg.Charset
}

func (pg *Page) GetStatusCode() int {
	return pg.Status
}

// AddTargetRequestUrls add raw url as a get request
func (pg *Page) AddTargetRequestUrls(urlStrs ...string) {
	for _, urlStr := range urlStrs {
		req, err := NewGetRequest(urlStr)
		if err != nil {
			log.Warnf("add target request url error: %s", err)
			continue
		}
		pg.AddTargetRequests(req)
	}
}

// AddTargetRequests add request
func (pg *Page) AddTargetRequests(requests ...IRequest) {
	for _, req := range requests {
		urlStr := req.GetUrl()
		if urlStr == "" || urlStr == "#" || strings.HasPrefix(urlStr, "javascript:") {
			continue
		}
		if urlStr = pg.canonicalizeUrl(urlStr); urlStr == "" {
			continue
		}
		pg.TargetRequests = append(pg.TargetRequests, req)
	}
}

func (pg *Page) canonicalizeUrl(url string) string {
	switch {
	case strings.Index(url, "http://") == 0 || strings.Index(url, "https://") == 0:
		return url
	case strings.Index(url, "/") == 0 || strings.Index(url, "?") == 0:
		return pg.domain + url
	default:
		return ""
	}
}

func (pg *Page) GetTargetRequests() []IRequest {
	return pg.TargetRequests
}

func (pg *Page) GetRespHeaders() map[string][]string {
	return pg.RespHeaders
}

func (pg *Page) GetUrl() *selector.PlainText {
	return pg.Url
}

func (pg *Page) StoreField(key string, obj interface{}) {
	pg.resultItems.Put(key, obj)
}

func (pg *Page) GetResultItems() IResultItems {
	return pg.resultItems
}

func (pg *Page) GetRawText() string {
	return pg.rawText
}

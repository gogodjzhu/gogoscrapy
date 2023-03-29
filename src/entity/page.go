package entity

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gogodjzhu/gogoscrapy/src/selector"
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

func (this *Page) GetRequest() IRequest {
	return this.Request
}

func (this *Page) GetDocument() *goquery.Document {
	return this.Document
}

func (this *Page) GetHtmlNode() *selector.HtmlNode {
	return this.HtmlRootNode
}

func (this *Page) GetCharset() string {
	return this.Charset
}

func (this *Page) GetStatusCode() int {
	return this.Status
}

//add raw request url
func (this *Page) AddTargetRequestUrls(urlStrs ...string) {
	for _, urlStr := range urlStrs {
		req := NewRequest(urlStr)
		this.AddTargetRequests(req)
	}
}

//add request
func (this *Page) AddTargetRequests(requests ...IRequest) {
	for _, req := range requests {
		urlStr := req.GetUrl()
		if urlStr == "" || urlStr == "#" || strings.HasPrefix(urlStr, "javascript:") {
			continue
		}
		if urlStr = this.canonicalizeUrl(urlStr); urlStr == "" {
			continue
		}
		this.TargetRequests = append(this.TargetRequests, req)
	}
}

func (this *Page) canonicalizeUrl(url string) string {
	switch {
	case strings.Index(url, "http://") == 0 || strings.Index(url, "https://") == 0:
		return url
	case strings.Index(url, "/") == 0 || strings.Index(url, "?") == 0:
		return this.domain + url
	default:
		return ""
	}
}

func (this *Page) GetTargetRequests() []IRequest {
	return this.TargetRequests
}

func (this *Page) GetRespHeaders() map[string][]string {
	return this.RespHeaders
}

func (this *Page) GetUrl() *selector.PlainText {
	return this.Url
}

func (this *Page) StoreField(key string, obj interface{}) {
	this.resultItems.Put(key, obj)
}

func (this *Page) GetResultItems() IResultItems {
	return this.resultItems
}

func (this *Page) GetRawText() string {
	return this.rawText
}

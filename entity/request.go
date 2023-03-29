package entity

import "net/http"

type IRequest interface {
	GetUrl() string
	GetMethod() string
	SetMethod(method string) IRequest
	GetExtras() map[string]interface{}
	SetExtras(extras map[string]interface{}) IRequest
	PutExtra(key string, value interface{}) IRequest
	GetCookies() map[string]string
	SetCookies(cookies map[string]string) IRequest
	PutCookie(key, value string) IRequest
	GetHeaders() map[string][]string
	SetHeaders(headers map[string][]string) IRequest
	PutHeader(key string, value []string) IRequest
	GetPriority() int64
	SetPriority(priority int64) IRequest
	GetCharset() string
	SetCharset(charset string) IRequest
	IsUseProxy() bool
	SetUseProxy(use bool) IRequest
	IsRetry() bool //this request is retry action
}

const CycleTriedTimes = "__cycle_tried_times"

type Request struct {
	Url      string
	Method   string
	Extras   map[string]interface{}
	Cookies  map[string]string
	Headers  map[string][]string
	Priority int64
	Charset  string
	UseProxy bool
}

func NewRequest(url string) IRequest {
	return &Request{
		Url:      url,
		Method:   http.MethodGet,
		Extras:   map[string]interface{}{},
		Cookies:  map[string]string{},
		Headers:  map[string][]string{},
		Priority: -1,
		Charset:  "UTF-8",
		UseProxy: false,
	}
}

func (this *Request) GetUrl() string {
	return this.Url
}

func (this *Request) GetMethod() string {
	return this.Method
}

func (this *Request) SetMethod(method string) IRequest {
	this.Method = method
	return this
}

func (this *Request) GetExtras() map[string]interface{} {
	return this.Extras
}

func (this *Request) SetExtras(extras map[string]interface{}) IRequest {
	this.Extras = extras
	return this
}

func (this *Request) PutExtra(key string, value interface{}) IRequest {
	this.Extras[key] = value
	return this
}

func (this *Request) GetCookies() map[string]string {
	return this.Cookies
}

func (this *Request) SetCookies(cookies map[string]string) IRequest {
	this.Cookies = cookies
	return this
}

func (this *Request) PutCookie(key, value string) IRequest {
	this.Cookies[key] = value
	return this
}

func (this *Request) GetHeaders() map[string][]string {
	return this.Headers
}

func (this *Request) SetHeaders(headers map[string][]string) IRequest {
	this.Headers = headers
	return this
}

func (this *Request) PutHeader(key string, value []string) IRequest {
	this.Headers[key] = value
	return this
}

func (this *Request) GetPriority() int64 {
	return this.Priority
}

func (this *Request) SetPriority(priority int64) IRequest {
	return this
}

func (this *Request) GetCharset() string {
	return this.Charset
}

func (this *Request) SetCharset(charset string) IRequest {
	this.Charset = charset
	return this
}

func (this *Request) IsUseProxy() bool {
	return this.UseProxy
}

func (this *Request) SetUseProxy(use bool) IRequest {
	this.UseProxy = use
	return this
}

func (this *Request) IsRetry() bool {
	return this.Extras[CycleTriedTimes] != nil
}

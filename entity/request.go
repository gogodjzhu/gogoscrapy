package entity

import (
	"errors"
	"github.com/gogodjzhu/gogoscrapy/utils"
	"net/http"
)

type IRequest interface {
	GetUrl() string
	GetMethod() string
	SetMethod(method string) IRequest
	GetBody() []byte
	SetBody(body []byte) IRequest
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
const AssignedProxy = "__assigned_proxy"

type Request struct {
	Url      string
	Method   string
	Body     []byte
	Extras   map[string]interface{}
	Cookies  map[string]string
	Headers  map[string][]string
	Priority int64
	Charset  string
	UseProxy bool
}

func NewGetRequest(url string) (IRequest, error) {
	if !utils.IsUrl(url) {
		return nil, errors.New("url is not valid: " + url)
	}
	return &Request{
		Url:      url,
		Method:   http.MethodGet,
		Body:     nil,
		Extras:   map[string]interface{}{},
		Cookies:  map[string]string{},
		Headers:  map[string][]string{},
		Priority: -1,
		Charset:  "UTF-8",
		UseProxy: false,
	}, nil
}

func (r *Request) GetUrl() string {
	return r.Url
}

func (r *Request) GetMethod() string {
	return r.Method
}

func (r *Request) SetMethod(method string) IRequest {
	r.Method = method
	return r
}

func (r *Request) GetBody() []byte {
	return r.Body
}

func (r *Request) SetBody(body []byte) IRequest {
	r.Body = body
	return r
}

func (r *Request) GetExtras() map[string]interface{} {
	return r.Extras
}

func (r *Request) SetExtras(extras map[string]interface{}) IRequest {
	r.Extras = extras
	return r
}

func (r *Request) PutExtra(key string, value interface{}) IRequest {
	r.Extras[key] = value
	return r
}

func (r *Request) GetCookies() map[string]string {
	return r.Cookies
}

func (r *Request) SetCookies(cookies map[string]string) IRequest {
	r.Cookies = cookies
	return r
}

func (r *Request) PutCookie(key, value string) IRequest {
	r.Cookies[key] = value
	return r
}

func (r *Request) GetHeaders() map[string][]string {
	return r.Headers
}

func (r *Request) SetHeaders(headers map[string][]string) IRequest {
	r.Headers = headers
	return r
}

func (r *Request) PutHeader(key string, value []string) IRequest {
	r.Headers[key] = value
	return r
}

func (r *Request) GetPriority() int64 {
	return r.Priority
}

func (r *Request) SetPriority(priority int64) IRequest {
	return r
}

func (r *Request) GetCharset() string {
	return r.Charset
}

func (r *Request) SetCharset(charset string) IRequest {
	r.Charset = charset
	return r
}

func (r *Request) IsUseProxy() bool {
	return r.UseProxy
}

func (r *Request) SetUseProxy(use bool) IRequest {
	r.UseProxy = use
	return r
}

func (r *Request) IsRetry() bool {
	return r.Extras[CycleTriedTimes] != nil
}

package entity

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gogodjzhu/gogoscrapy/selector"
	"testing"
)

func TestPage_canonicalizeUrl(t *testing.T) {
	type fields struct {
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
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "ok-http",
			fields: fields{},
			args: args{
				url: "http://gogodjzhu.com",
			},
			want: "http://gogodjzhu.com",
		},
		{
			name:   "ok-https",
			fields: fields{},
			args: args{
				url: "https://gogodjzhu.com",
			},
			want: "https://gogodjzhu.com",
		},
		{
			name: "ok-relative",
			fields: fields{
				domain: "https://gogodjzhu.com",
			},
			args: args{
				url: "/gogodjzhu",
			},
			want: "https://gogodjzhu.com/gogodjzhu",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := &Page{
				Request:        tt.fields.Request,
				Document:       tt.fields.Document,
				HtmlRootNode:   tt.fields.HtmlRootNode,
				Charset:        tt.fields.Charset,
				Status:         tt.fields.Status,
				TargetRequests: tt.fields.TargetRequests,
				RespHeaders:    tt.fields.RespHeaders,
				Url:            tt.fields.Url,
				rawText:        tt.fields.rawText,
				resultItems:    tt.fields.resultItems,
				domain:         tt.fields.domain,
			}
			if got := pg.canonicalizeUrl(tt.args.url); got != tt.want {
				t.Errorf("canonicalizeUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

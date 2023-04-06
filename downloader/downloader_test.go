package downloader

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gogodjzhu/gogoscrapy/entity"
	"net/http"
	"strings"
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

func Test_getRawText(t *testing.T) {
	type args struct {
		header http.Header
		doc    *goquery.Document
	}
	jsonDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(`{"name":"gogodjzhu"}`))
	htmlDoc, _ := goquery.NewDocumentFromReader(strings.NewReader("<html><head><title>test</title></head><body><p>gogodjzhu</p></body></html>"))
	xmlDoc, _ := goquery.NewDocumentFromReader(strings.NewReader("<xml><name>gogodjzhu</name></xml>"))
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok-json",
			args: args{
				header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				doc: jsonDoc,
			},
			want:    `{"name":"gogodjzhu"}`,
			wantErr: false,
		},
		{
			name: "ok-html",
			args: args{
				header: http.Header{
					"Content-Type": []string{"text/html"},
				},
				doc: htmlDoc,
			},
			want:    "<p>gogodjzhu</p>",
			wantErr: false,
		},
		{
			name: "ok-xml",
			args: args{
				header: http.Header{
					"Content-Type": []string{"text/xml"},
				},
				doc: xmlDoc,
			},
			want:    "<xml><name>gogodjzhu</name></xml>",
			wantErr: false,
		},
		{
			name: "default",
			args: args{
				header: http.Header{},
				doc:    htmlDoc,
			},
			want:    "<html><head><title>test</title></head><body><p>gogodjzhu</p></body></html>",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRawText(tt.args.header, tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRawText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getRawText() got = %v, want %v", got, tt.want)
			}
		})
	}
}

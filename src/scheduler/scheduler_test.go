package scheduler

import (
	"github.com/gogodjzhu/gogoscrapy/src/entity"
	"net/http"
	"testing"
)

func TestQueueScheduler_PushWhenNoDuplicate(t *testing.T) {
	queueScheduler := NewQueueScheduler()
	queueScheduler.Push(&entity.Request{Url: "http://gogodjzhu.com"})                          //success
	queueScheduler.Push(&entity.Request{Url: "http://gogodjzhu.com", Method: http.MethodPost}) //success for post
	queueScheduler.Push(&entity.Request{Url: "http://gogodjzhu.com"})                          //duplicated
	req := queueScheduler.Poll()
	if req.GetUrl() != "http://gogodjzhu.com" {
		t.Error("test failed @ TestQueueScheduler_PushWhenNoDuplicate")
	}
	req = queueScheduler.Poll()
	if req.GetUrl() != "http://gogodjzhu.com" {
		t.Error("test failed @ TestQueueScheduler_PushWhenNoDuplicate")
	}
	req = queueScheduler.Poll()
	if req != nil {
		t.Error("test failed @ TestQueueScheduler_PushWhenNoDuplicate")
	}
}

package scheduler

import (
	"github.com/gogodjzhu/gogoscrapy/entity"
	"testing"
)

func TestNewMemDuplicateRemover(t *testing.T) {
	remover := NewMemDuplicateRemover()
	if result, err := remover.IsDuplicate(&entity.Request{Url: "abc"}); result || err != nil {
		t.FailNow()
	}
	if result, err := remover.IsDuplicate(&entity.Request{Url: "abc"}); !result || err != nil {
		t.FailNow()
	}
	if result, err := remover.GetTotalCount(); result != 1 || err != nil {
		t.FailNow()
	}
	if err := remover.ResetDuplicate(); err != nil {
		t.FailNow()
	}
	if result, err := remover.IsDuplicate(&entity.Request{Url: "abc"}); result || err != nil {
		t.FailNow()
	}
}

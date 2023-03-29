package scheduler

import (
	"github.com/gogodjzhu/gogoscrapy/src/entity"
	"testing"
)

func TestNewMemDuplicateRemover(t *testing.T) {
	remover := NewMemDuplicateRemover()
	if remover.IsDuplicate(&entity.Request{Url: "abc"}) {
		t.Error("test failed @ TestNewMemDuplicateRemover")
	}
	if !remover.IsDuplicate(&entity.Request{Url: "abc"}) {
		t.Error("test failed @ TestNewMemDuplicateRemover")
	}
	if remover.GetTotalCount() != 1 {
		t.Error("test failed @ TestNewMemDuplicateRemover")
	}
	remover.ResetDuplicate()
	if remover.IsDuplicate(&entity.Request{Url: "abc"}) {
		t.Error("test failed @ TestNewMemDuplicateRemover")
	}
}

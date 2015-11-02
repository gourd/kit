package context_test

import (
	"testing"
	"time"

	gcontext "github.com/gourd/kit/context"
)

func TestEmptyContextDeadline(t *testing.T) {
	ctx := gcontext.New(nil)
	dl, ok := ctx.Deadline()
	if ok != false {
		t.Error("default ok (in `_, ok := context.Deadline()`) should be false")
	}
	if !dl.Equal(time.Unix(0, 0)) {
		t.Error("default deadline time should be unix epoch")
	}
}

func TestEmptyContextDone(t *testing.T) {
	ctx := gcontext.New(nil)
	done := ctx.Done()
	if done != nil {
		t.Error("ctx.Done expected to return nil, got %#v", done)
	}
}

func TestEmptyContextErr(t *testing.T) {
	ctx := gcontext.New(nil)
	err := ctx.Err()
	if err != nil {
		t.Error("ctx.Err expected to return nil, got %#v", err)
	}
}

func TestEmptyContextValue(t *testing.T) {
	ctx := gcontext.New(nil)
	val := ctx.Value("anything")
	if val != nil {
		t.Error("ctx.Value(\"anything\") expected to return nil, got %#v", val)
	}
}

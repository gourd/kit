package gourdctx_test

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
	if want, have := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), dl; !have.Equal(want) {
		t.Errorf("unexpected deadline time\nExpect: %s\nGot:    %s", want, have)
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

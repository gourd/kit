package upperio_test

import (
	"github.com/gourd/kit/store"
	"github.com/gourd/kit/store/upperio"

	"testing"
)

func TestSort(t *testing.T) {
	q := store.NewQuery()
	q.GetSorts().Add("para1").Add("-para2").Add("-para3")

	res := upperio.Sort(q)
	expLen := 3
	if l := len(res); l != expLen {
		t.Errorf("result length expected: %d, get: %d", expLen, l)
		t.FailNow()
	}

	if expStr := "para1"; res[0] != expStr {
		t.Errorf("result[0] expected: %d, get: %d", expStr, res[0])
	}
	if expStr := "-para2"; res[1] != expStr {
		t.Errorf("result[1] expected: %d, get: %d", expStr, res[1])
	}
	if expStr := "-para3"; res[2] != expStr {
		t.Errorf("result[2] expected: %d, get: %d", expStr, res[2])
	}
}

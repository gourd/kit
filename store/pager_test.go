package store_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/gourd/kit/store"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestPager_Basic(t *testing.T) {
	total := rand.Intn(10000) + 1
	offset := rand.Intn(10000) + 1
	limit := rand.Intn(10000) + 1

	p1 := store.NewPager()
	p1.SetTotal(total)
	p1.SetLimit(limit, offset)
	data, err := json.Marshal(p1)
	if err != nil {
		t.Errorf("marshal error: %#v", err.Error())
	}

	p2 := store.NewPager()
	err = json.Unmarshal(data, &p2)
	if err != nil {
		t.Errorf("unmarshal error: %#v", err.Error())
	}
	if want, have := total, p2.GetTotal(); want != have {
		t.Errorf("total want: %#v, got: %#v", want, have)
	}
	limit2, offset2 := p2.GetLimit()
	if want, have := limit, limit2; want != have {
		t.Errorf("limit want: %#v, got: %#v", want, have)
	}
	if want, have := offset, offset2; want != have {
		t.Errorf("offset want: %#v, got: %#v", want, have)
	}
}

func testPagerHasOnly(p store.Pager, keys ...string) error {
	vmap := make(map[string]interface{})
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshal error: %#v", err.Error())
	}

	if err = json.Unmarshal(data, &vmap); err != nil {
		return fmt.Errorf("unmarshal error: %#v", err.Error())
	} else if want, have := len(keys), len(vmap); want != have {
		return fmt.Errorf("want: %#v, got: %#v", want, have)
	} else {
		for _, key := range keys {
			if _, ok := vmap[key]; !ok {
				return fmt.Errorf("key %#v not exists. (vmap = %#v)",
					key, vmap)
			}
		}
	}
	return nil
}

func TestPager_Partial(t *testing.T) {

	p1 := store.NewPager().SetTotal(rand.Intn(10000) + 1)
	if err := testPagerHasOnly(p1, "total"); err != nil {
		t.Error(err)
	}

	p2 := store.NewPager().SetLimit(
		rand.Intn(10000)+1, rand.Intn(10000)+1)
	if err := testPagerHasOnly(p2, "limit", "offset"); err != nil {
		t.Error(err)
	}

}

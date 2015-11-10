package store_test

import (
	"github.com/gourd/kit/store"

	"encoding/json"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestResponseMarshaler(t *testing.T) {
	var m json.Marshaler = store.NewResponse("", nil)
	_ = m
	t.Log("*Response implements json.Marshaler")
}

func TestExpandResponse(t *testing.T) {
	listIn := []int{rand.Intn(100), rand.Intn(100), rand.Intn(100), rand.Intn(100)}
	vmap := map[string]interface{}{
		"list": listIn,
	}
	sres := store.ExpandResponse(vmap)

	// test unset key
	if want, have := interface{}(nil), sres.Get("nothing"); want != have {
		t.Errorf("want: %#v, got: %#v", want, have)
	}

	// test getting the list back
	if listGet, ok := sres.Get("list").([]int); !ok {
		t.Errorf("want: []int, got: %#v", sres.Get("list"))
	} else if want, have := len(listIn), len(listGet); want != have {
		t.Errorf("list lenght wanted: %#v, got: %#v", want, have)
	} else {
		for i, item := range listIn {
			if want, have := item, listGet[i]; want != have {
				t.Errorf("list[%d] want: %#v, got: %#v", i, want, have)
			}
		}
	}

	// test marshal to map[string]interface{}
	data, err := json.Marshal(&sres)
	if err != nil {
		t.Errorf("Marshal error %#v", err.Error())
	}
	v := make(map[string]interface{})
	json.Unmarshal(data, &v)

	// test status
	if status, ok := v["status"]; !ok {
		t.Error("status not found")
	} else if floatStatus, ok := status.(float64); !ok {
		t.Error("status is not number")
	} else if want, have := http.StatusOK, int(floatStatus); want != have {
		t.Errorf("want status %#v, got %#v", want, have)
	}

	// test list content
	if resList, ok := v["list"]; !ok {
		t.Error("list not found")
	} else if listOut, ok := resList.([]interface{}); !ok {
		t.Errorf("list want []interface{}, got %#v", resList)
	} else if want, have := len(listIn), len(listOut); want != have {
		t.Errorf("list lenght wanted: %#v, got: %#v", want, have)
	} else {
		for i, item := range listIn {
			if v, ok := listOut[i].(float64); !ok {
				t.Errorf("list[%d] want float64, got: %#v", listOut[i])
			} else if want, have := item, int(v); want != have {
				t.Errorf("list[%d] want: %#v, got: %#v", i, want, have)
			}
		}
	}

}

func TestNewResponse(t *testing.T) {
	listIn := []int{rand.Intn(100), rand.Intn(100), rand.Intn(100), rand.Intn(100)}
	sres := store.NewResponse("list", listIn)

	// test unset key
	if want, have := interface{}(nil), sres.Get("nothing"); want != have {
		t.Errorf("want: %#v, got: %#v", want, have)
	}

	// test getting the list back
	if listGet, ok := sres.Get("list").([]int); !ok {
		t.Errorf("want: []int, got: %#v", sres.Get("list"))
	} else if want, have := len(listIn), len(listGet); want != have {
		t.Errorf("list lenght wanted: %#v, got: %#v", want, have)
	} else {
		for i, item := range listIn {
			if want, have := item, listGet[i]; want != have {
				t.Errorf("list[%d] want: %#v, got: %#v", i, want, have)
			}
		}
	}

	// test marshal to map[string]interface{}
	data, err := json.Marshal(&sres)
	if err != nil {
		t.Errorf("Marshal error %#v", err.Error())
	}
	v := make(map[string]interface{})
	json.Unmarshal(data, &v)

	// test status
	if status, ok := v["status"]; !ok {
		t.Error("status not found")
	} else if floatStatus, ok := status.(float64); !ok {
		t.Error("status is not number")
	} else if want, have := http.StatusOK, int(floatStatus); want != have {
		t.Errorf("want status %#v, got %#v", want, have)
	}

	// test list content
	if resList, ok := v["list"]; !ok {
		t.Error("list not found")
	} else if listOut, ok := resList.([]interface{}); !ok {
		t.Errorf("list want []interface{}, got %#v", resList)
	} else if want, have := len(listIn), len(listOut); want != have {
		t.Errorf("list lenght wanted: %#v, got: %#v", want, have)
	} else {
		for i, item := range listIn {
			if v, ok := listOut[i].(float64); !ok {
				t.Errorf("list[%d] want float64, got: %#v", listOut[i])
			} else if want, have := item, int(v); want != have {
				t.Errorf("list[%d] want: %#v, got: %#v", i, want, have)
			}
		}
	}

}

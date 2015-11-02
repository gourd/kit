package store_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gourd/kit/store"
)

func TestError(t *testing.T) {
	estatus := rand.Intn(900)
	ecode := estatus*100 + rand.Intn(100)
	emsg := fmt.Sprintf("random error %d", ecode)

	var err error = store.Error(ecode, emsg)
	code, msg := store.ParseError(err)

	if code != ecode {
		t.Errorf("code output not correct. Expect %#v but get %#v",
			ecode, code)
	}
	if msg != emsg {
		t.Errorf("msg output not correct. Expect %#v but get %#v",
			emsg, msg)
	}

	serr := store.ExpandError(err)
	if serr.Status != estatus {
		t.Errorf("Incorrect StoreError.Status. Expecting %#v but get %#v",
			estatus, serr.Status)
	}
	if serr.Code != ecode {
		t.Errorf("Incorrect StoreError.Code. Expecting %#v but get %#v",
			ecode, serr.Code)
	}
	if serr.ServerMsg != emsg {
		t.Errorf("Incorrect StoreError.ServerMsg. Expecting %#v but get %#v",
			emsg, serr.ServerMsg)
	}
	if serr.ClientMsg != emsg {
		t.Errorf("Incorrect StoreError.ClientMsg. Expecting %#v but get %#v",
			emsg, serr.ClientMsg)
	}
	if serr.DeveloperMsg != "" {
		t.Errorf("Incorrect StoreError.DeveloperMsg. Expecting %#v but get %#v",
			"", serr.DeveloperMsg)
	}

}

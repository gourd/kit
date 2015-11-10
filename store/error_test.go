package store_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gourd/kit/store"
)

func TestError(t *testing.T) {
	estatus := (rand.Intn(9) * 100) + (rand.Intn(9) * 10)
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

func TestErrorFormat(t *testing.T) {
	tpl := "(%s) some custom bad request: %#v"
	msg := "hello world"

	estatus := rand.Intn(900)
	ecode := estatus*100 + rand.Intn(100)

	err := store.Error(ecode, "").
		TellClient(tpl, "client", msg).
		TellDeveloper(tpl, "developer", msg).
		TellServer(tpl, "server", msg)

	if want, have := fmt.Sprintf(tpl, "client", msg), err.String(); want != have {
		t.Errorf("want: %#v, got %#v", want, have)
	}

	if want, have := fmt.Sprintf(tpl, "client", msg), fmt.Sprintf("%s", err); want != have {
		t.Errorf("want: %#v, got %#v", want, have)
	}

	if want, have := fmt.Sprintf(tpl, "client", msg), err.Error(); want != have {
		t.Errorf("want: %#v, got %#v", want, have)
	}

	if want, have := fmt.Sprintf(tpl, "client", msg), err.ClientMsg; want != have {
		t.Errorf("want: %#v, got %#v", want, have)
	}

	if want, have := fmt.Sprintf(tpl, "developer", msg), err.DeveloperMsg; want != have {
		t.Errorf("want: %#v, got %#v", want, have)
	}

	if want, have := fmt.Sprintf(tpl, "server", msg), err.ServerMsg; want != have {
		t.Errorf("want: %#v, got %#v", want, have)
	}
}

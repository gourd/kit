package store_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gourd/kit/store"
)

func TestErrorSingletons(t *testing.T) {
	a := store.ErrorNotFound
	b := func(in error) (out error) {
		out = in
		return out
	}(a)
	c := store.ErrorNotFound

	if a != b {
		t.Errorf("Singleton failed to be passed to another error variable")
	}
	if a != c {
		t.Errorf("Singleton failed to equal to itself")
	}
}

func TestParseError_Singletons(t *testing.T) {

	var code int
	var msg string

	var err error = store.ErrorNotFound
	msgExpect := err.Error()

	code, msg = store.ParseError(err)
	if code != http.StatusNotFound {
		t.Errorf("Incorrect status code. Expecting %d but get %d",
			http.StatusNotFound, code)
	}
	if msg != msgExpect {
		t.Errorf("Incorrect status message. Expecting %s but get %s",
			msgExpect, msg)
	}

	serr := store.ExpandError(err)
	if serr.Status != http.StatusNotFound {
		t.Errorf("Incorrect StoreError.Status. Expecting %#v but get %#v",
			http.StatusNotFound, serr.Status)
	}
	if serr.Code != http.StatusNotFound {
		t.Errorf("Incorrect StoreError.Code. Expecting %#v but get %#v",
			http.StatusNotFound, serr.Code)
	}
	if serr.ServerMsg != msgExpect {
		t.Errorf("Incorrect StoreError.ServerMsg. Expecting %#v but get %#v",
			msgExpect, serr.ServerMsg)
	}
	if serr.ClientMsg != msgExpect {
		t.Errorf("Incorrect StoreError.ClientMsg. Expecting %#v but get %#v",
			msgExpect, serr.ClientMsg)
	}
	if serr.DeveloperMsg != "" {
		t.Errorf("Incorrect StoreError.DeveloperMsg. Expecting %#v but get %#v",
			"", serr.DeveloperMsg)
	}

}

func TestParseError_User(t *testing.T) {

	var code int
	var msg string
	userMsg := "Some user error"

	var err error = errors.New(userMsg)
	code, msg = store.ParseError(err)

	if code != http.StatusInternalServerError {
		t.Errorf("Incorrect status code. Expecting %d but get %d",
			http.StatusInternalServerError, code)
	}
	if msg != userMsg {
		t.Errorf("Incorrect status message. Expecting %s but get %s",
			userMsg, msg)
	}

	serr := store.ExpandError(err)
	if serr.Status != http.StatusInternalServerError {
		t.Errorf("Incorrect StoreError.Status. Expecting %#v but get %#v",
			http.StatusInternalServerError, serr.Status)
	}
	if serr.Code != http.StatusInternalServerError {
		t.Errorf("Incorrect StoreError.Code. Expecting %#v but get %#v",
			http.StatusInternalServerError, serr.Code)
	}
	if serr.ServerMsg != userMsg {
		t.Errorf("Incorrect StoreError.ServerMsg. Expecting %#v but get %#v",
			userMsg, serr.ServerMsg)
	}
	if serr.ClientMsg != userMsg {
		t.Errorf("Incorrect StoreError.ClientMsg. Expecting %#v but get %#v",
			userMsg, serr.ClientMsg)
	}
	if serr.DeveloperMsg != "" {
		t.Errorf("Incorrect StoreError.DeveloperMsg. Expecting %#v but get %#v",
			"", serr.DeveloperMsg)
	}

}

func TestParseError_nil(t *testing.T) {
	var code int
	var msg string

	err := store.StatusFound
	msgExpect := err.Error()

	code, msg = store.ParseError(nil)

	if code != http.StatusFound {
		t.Errorf("Incorrect status code. Expecting %d but get %d",
			http.StatusFound, code)
	}
	if msg != msgExpect {
		t.Errorf("Incorrect status message. Expecting %s but get %s",
			msgExpect, msg)
	}

}

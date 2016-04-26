package httpservice_test

import (
	"testing"
	"time"

	"github.com/gourd/kit/service/http"
)

func TestEnforceCreate(t *testing.T) {
	type myType struct {
		ID        string    `json:"id,omitempty"`
		Name      string    `json:"name"`
		Created   time.Time `json:"created" gourdcreate:"now"`
		Updated   time.Time `json:"updated" gourdcreate:"now"`
		Published time.Time `json:"created" gourdcreate:"now,omitnotempty"`
	}

	now := time.Now()

	// test with empty Published
	payload1 := &myType{}
	if err := httpservice.EnforceCreate(payload1); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if want, have := now.Unix(), payload1.Created.Unix(); want != have {
		t.Errorf("Created: expected %#v, got %#v", want, have)
	}
	if want, have := now.Unix(), payload1.Updated.Unix(); want != have {
		t.Errorf("Updated: expected %#v, got %#v", want, have)
	}
	if want, have := now.Unix(), payload1.Published.Unix(); want != have {
		t.Errorf("Published: expected %#v, got %#v", want, have)
	}

	// test with non-empty Published
	payload2 := &myType{Published: time.Unix(0, 0)}
	if err := httpservice.EnforceCreate(payload2); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if want, have := now.Unix(), payload2.Created.Unix(); want != have {
		t.Errorf("Created: expected %#v, got %#v", want, have)
	}
	if want, have := now.Unix(), payload2.Updated.Unix(); want != have {
		t.Errorf("Updated: expected %#v, got %#v", want, have)
	}
	if want, have := time.Unix(0, 0).Unix(), payload2.Published.Unix(); want != have {
		t.Errorf("Published: expected %#v, got %#v", want, have)
	}
}

func TestEnforceUpdate(t *testing.T) {
	type myType struct {
		ID        string    `json:"id,omitempty" gourdupdate:"preserve"`
		Name      string    `json:"name"`
		Created   time.Time `json:"created" gourdupdate:"preserve"`
		Updated   time.Time `json:"updated" gourdupdate:"now"`
		Published time.Time `json:"created"`
	}

	type myType2 struct {
		ID string `json:"id,omitempty" gourdupdate:"preserve"`
	}

	now := time.Now()
	str := "hello"

	if want, have := "reflect.ValueOf(update) is of zero Value", httpservice.EnforceUpdate(str, nil).Error(); want != have {
		t.Errorf("expected error message %#v, got %#v", want, have)
	}
	if want, have := "reflect.ValueOf(update) is of zero Value", httpservice.EnforceUpdate(&str, nil).Error(); want != have {
		t.Errorf("expected error message %#v, got %#v", want, have)
	}

	v1, v2 := &myType{}, &myType2{}
	if want, have := "*original (httpservice_test.myType) is not of same type of *update (httpservice_test.myType2)", httpservice.EnforceUpdate(v1, v2).Error(); want != have {
		t.Errorf("expected error message %#v, got %#v", want, have)
	}

	// test with empty Published
	original := &myType{
		ID:        "abcd123",
		Created:   time.Unix(1024, 0),
		Updated:   time.Unix(1024, 0),
		Published: time.Unix(1024, 0),
	}
	update := &myType{
		ID:        "abcd1245",
		Created:   time.Unix(2048, 0),
		Published: time.Unix(1024, 0),
	}

	if err := httpservice.EnforceUpdate(original, update); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if want, have := original.ID, update.ID; want != have {
		t.Errorf("ID: expected %#v, got %#v", want, have)
	}
	if want, have := original.Created, update.Created; want != have {
		t.Errorf("Created: expected %s, got %s", want, have)
	}
	if want, have := now.Unix(), update.Updated.Unix(); want != have {
		t.Errorf("Updated: expected %#v, got %#v", want, have)
	}
	if want, have := original.Published, update.Published; want != have {
		t.Errorf("Published: expected %s, got %s", want, have)
	}
}

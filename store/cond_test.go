package store_test

import (
	"github.com/gourd/kit/store"

	"testing"
)

func TestBasicConds(t *testing.T) {
	t.Parallel()
	var c store.Conds
	c = &store.BasicConds{}
	t.Logf("BasicConds can be casted to Conds: %#v", c)
}

func TestNewConds(t *testing.T) {
	t.Parallel()
	var c store.Conds
	c = store.NewConds()
	t.Logf("NewCond can return Conds: %#v", c)
}

func TestBasicConds_AddGetAll(t *testing.T) {
	t.Parallel()
	c := store.NewConds().Add("foo", "bar").Add("hello", "world")
	a := c.GetAll()

	if a[0].Prop != "foo" {
		t.Errorf("Failed to add: %#v", a[0])
	} else if a[0].Value != "bar" {
		t.Errorf("Failed testing value with original string: %#v -> %s", a[0].Value, a[0].Value)
	}

	if a[1].Prop != "hello" {
		t.Errorf("Failed to add: %#v", a[1])
	} else if a[1].Value != "world" {
		t.Errorf("Failed testing value with original string: %#v -> %s", a[1].Value, a[1].Value)
	}
}

func TestBasicConds_AddGetMap(t *testing.T) {
	t.Parallel()
	c := store.NewConds().Add("foo", "bar").Add("hello", "world")
	m, err := c.GetMap()

	if err != nil {
		t.Errorf("Error in GetMap(): %s", err.Error())
	}

	if m["foo"] != "bar" {
		t.Errorf("Failed to get proper map: m[\"foo\"] is \"%#v\" instead of \"%s\"",
			m["foo"])
	}
	if m["hello"] != "world" {
		t.Errorf("Failed to get proper map: m[\"hello\"] is \"%#v\" instead of \"%s\"",
			m["hello"])
	}
}

func TestBasicConds_AddGetMapErr(t *testing.T) {
	t.Parallel()
	c := store.NewConds().Add("foo", "bar").Add("foo", "again")
	_, err := c.GetMap()

	if err == nil {
		t.Errorf("Failed to return error with conflicting map conditions")
	}
}

func TestBasicConds_SetGetRel(t *testing.T) {
	t.Parallel()
	c := store.NewConds().Add("foo", "bar").Add("hello", "world")
	if c.GetRel() != store.And {
		t.Errorf("Conds Rel flag is not initialized as And")
	} else {
		t.Log("Conds Rel initialized as And")
	}

	c.SetRel(store.Or)
	if c.GetRel() != store.Or {
		t.Errorf("Failed to set Conds Rel to Or")
	} else {
		t.Log("Conds Rel changed to Or")
	}
}

package store_test

import (
	"github.com/gourd/kit/store"

	"testing"
)

func TestBasicQuery_Limit(t *testing.T) {
	q := &store.BasicQuery{}
	success := true

	if q.Limit != 0 {
		t.Error("Limit is not initialized as 0")
		success = false
	} else if q.GetLimit() != 0 {
		t.Error("GetLimit did not get correct value")
		success = false
	}

	q.SetLimit(123)
	if q.Limit != 123 {
		t.Error("Limit is not set as 123")
		success = false
	} else if q.GetLimit() != 123 {
		t.Error("GetLimit did not get correct value")
		success = false
	}

	if success {
		t.Log("*BasicQuery Limit setter and getter works")
	}
}

func TestBasicQuery_Offset(t *testing.T) {
	q := &store.BasicQuery{}
	success := true

	if q.Offset != 0 {
		t.Error("Offset is not initialized as 0")
		success = false
	} else if q.GetOffset() != 0 {
		t.Error("GetOffset did not get correct value")
		success = false
	}

	q.SetOffset(123)
	if q.Offset != 123 {
		t.Error("Offset is not set as 123")
		success = false
	} else if q.GetOffset() != 123 {
		t.Error("GetOffset did not get correct value")
		success = false
	}

	if success {
		t.Log("*BasicQuery Offset setter and getter works")
	}
}

func TestBasicQuery_Conds(t *testing.T) {
	q := store.NewQuery().
		AddCond("HelloWorld =", "foo bar").
		AddCond("FooBar !=", "hello world")
	cs := q.GetConds().GetAll()
	success := true

	if cs[0].Prop != "HelloWorld =" {
		t.Errorf("First cond Prop is not defined correctly. "+
			"Expected \"%s\" but get \"%s\"",
			"HelloWorld =", cs[0].Prop)
		success = false
	} else if cs[0].Value != "foo bar" {
		t.Errorf("First cond Value is not defined correctly. "+
			"Expected \"%s\" but get \"%s\"",
			"foobar", cs[0].Value)
		success = false
	}

	if cs[1].Prop != "FooBar !=" {
		t.Errorf("Second cond Prop is not defined correctly. "+
			"Expected \"%s\" but get \"%s\"",
			"FooBar !=", cs[1].Prop)
		success = false
	} else if cs[1].Value != "hello world" {
		t.Errorf("Second cond Value is not defined correctly. "+
			"Expected \"%s\" but get \"%s\"",
			"hello world", cs[1].Value)
		success = false
	}

	if success {
		t.Log("Query cond routine works expectedly")
	}
}

func TestBasicQuery_Sorts(t *testing.T) {
	q := store.NewQuery().
		Sort("HelloWorld").
		Sort("-FooBar")
	ss := q.GetSorts().GetAll()
	success := true

	if ss[0].String() != "HelloWorld" {
		t.Errorf("First sort is not defined correctly. "+
			"Expected \"%s\" but get \"%s\"",
			"HelloWorld", ss[0])
		success = false
	}

	if ss[1].String() != "-FooBar" {
		t.Errorf("Second cond Prop is not defined correctly. "+
			"Expected \"%s\" but get \"%s\"",
			"-FooBar", ss[1])
		success = false
	}

	if success {
		t.Log("Query cond routine works expectedly")
	}
}

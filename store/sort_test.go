package store_test

import (
	"github.com/gourd/kit/store"

	"fmt"
	"testing"
)

func TestSort_AscDesc(t *testing.T) {
	s := &store.Sort{}
	success := true

	if store.Asc == store.Desc {
		t.Error("constant Asc = Desc")
		success = false
	}

	if s.Order != store.Asc {
		t.Error("Sort.Order is not initialized as Asc")
		success = false
	}

	s.Desc()
	if s.Order != store.Desc {
		t.Error("Sort.Order is not set as Desc")
		success = false
	}

	s.Asc()
	if s.Order != store.Asc {
		t.Error("Sort.Order is not set as Asc")
		success = false
	}

	if success {
		t.Log("Asc and Desc methods work expectedly")
	}
}

func TestSort_String(t *testing.T) {
	s := store.Sort{
		Name: "HelloWorld",
	}
	success := true

	if s.String() != "HelloWorld" {
		t.Error("Sorting asc is not displayed correctly by Sort.String()")
		success = false
	}

	s.Desc()
	if s.String() != "-HelloWorld" {
		t.Error("Sorting desc is not displayed correctly by Sort.String()")
		success = false
	}

	if success {
		t.Log("String methods work expectedly")
	}
}

func TestSortBy(t *testing.T) {
	var s *store.Sort = store.SortBy("Hello")
	_ = s
	t.Log("SortBy works expectedly")
}

func TestBasicSorts(t *testing.T) {
	var ss store.Sorts = &store.BasicSorts{}
	_ = ss
	t.Log("*BasicSorts implements Sorts")
}

func TestSorts_Routine(t *testing.T) {
	ss := &store.BasicSorts{}
	ss.
		Add("HelloWorld").
		Add("-FooBar")
	sli := ss.GetAll()
	if sli[0].String() != "HelloWorld" {
		t.Error("Incorrect first sort argument")
	} else if sli[1].String() != "-FooBar" {
		t.Error("Incorrect second sort argument")
	} else if fmt.Sprintf("%s", sli[1]) != "-FooBar" {
		t.Error("Sort failed to be formated as string")
	} else {
		t.Log("routine works")
	}
}

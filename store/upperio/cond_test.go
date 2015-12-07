package upperio_test

import (
	"github.com/gourd/kit/store"
	"github.com/gourd/kit/store/upperio"
	"upper.io/db"

	"os"
	"testing"
)

func TestConds_empty(t *testing.T) {
	q := store.NewQuery()
	if cs := upperio.Conds(q.GetConds()); cs != nil {
		t.Errorf("Conds with empty new query should return nil. Instead got %#v", cs)
	}
}

func TestConds(t *testing.T) {

	var err error

	fn := "./test2.tmp"

	q := store.NewQuery().
		AddCond("HelloWorld =", "foo bar").
		AddCond("FooBar !=", "hello world")

	// dummy source
	source := upperio.Source(testUpperDb(fn))

	// add dummy data to the database
	if err := testUpperDbData(source); err != nil {
		t.Fatal(err.Error())
	}

	// connect to database again
	conn, err := source()
	if err != nil {
		t.Error(err.Error())
	}
	defer conn.Close()

	// query connection
	sess := conn.Raw().(db.Database)
	coll, err := sess.Collection("dummy_data")
	if err != nil {
		t.Error(err.Error())
	}

	conds := upperio.Conds(q.GetConds())
	res := coll.Find(conds)
	var tds []testData
	res.All(&tds)

	if len(tds) != 1 {
		t.Errorf("Incorrect test data set: %#v", tds)
	}

	// clean up the temp database
	err = os.Remove(fn)
	if err != nil {
		t.Error(err.Error())
	}

}

func TestConds_branching(t *testing.T) {

	var err error

	fn := "./test3.tmp"

	// two branch query
	cond1 := store.NewConds().
		Add("HelloWorld =", "foo bar").
		Add("FooBar !=", "hello world")
	cond2 := store.NewConds().
		Add("HelloWorld =", "foo bar 2").
		Add("FooBar !=", "hello world")

	q := store.NewQuery().
		AddCond("", cond1).
		AddCond("", cond2)

	q.GetConds().SetRel(store.Or)

	// test source
	source := upperio.Source(testUpperDb(fn))

	// add dummy data to the database
	if err := testUpperDbData(source); err != nil {
		t.Fatal(err.Error())
	}

	// connect to database again
	conn, err := source()
	if err != nil {
		t.Error(err.Error())
	}
	defer conn.Close()

	// query connection
	sess := conn.Raw().(db.Database)
	coll, err := sess.Collection("dummy_data")
	res := coll.Find(upperio.Conds(q.GetConds()))
	var tds []testData
	res.All(&tds)

	expLen := 2
	if l := len(tds); l != expLen {
		t.Errorf("result set size expected: %d, got: %d\ntest data set:\t%#v",
			expLen, l, tds)
	}

	// clean up the temp database
	err = os.Remove(fn)
	if err != nil {
		t.Error(err.Error())
	}

}

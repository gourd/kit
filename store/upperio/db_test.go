package upperio_test

import (
	"github.com/gourd/kit/store/upperio"

	"database/sql"
	"math/rand"
	"net/http"
	"testing"
	"time"
	"upper.io/db/sqlite"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func testUpperDb(fn string) (srcName string) {

	// testing database name
	srcName = RandStringRunes(10)

	// define a database source
	upperio.Define(
		srcName,
		sqlite.Adapter,
		sqlite.ConnectionURL{
			Database: fn,
		},
	)

	return
}

func testUpperDbData(t *testing.T, srcName string) string {

	// dummy request
	r := &http.Request{}

	sess, err := upperio.Open(r, srcName)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer upperio.Close(r, srcName)

	drv := sess.Driver().(*sql.DB)
	_, err = drv.Exec(`
		CREATE TABLE dummy_data (
			HelloWorld text,
			FooBar text,
			Data text
		)
	`)
	if err != nil {
		t.Fatalf(err.Error())
	}

	coll, err := sess.Collection("dummy_data")
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = coll.Append(&testData{
		HelloWorld: "foo bar",
		FooBar:     "foo bar",
		Data:       "something",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = coll.Append(&testData{
		HelloWorld: "foo bar 2",
		FooBar:     "foo bar 2",
		Data:       "something 2",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = coll.Append(&testData{
		HelloWorld: "foo bar 3",
		FooBar:     "foo bar 3",
		Data:       "something 3",
	})
	if err != nil {
		t.Fatalf(err.Error())
	}

	return srcName
}

// for dummy database content
type testData struct {
	HelloWorld string
	FooBar     string
	Data       string
}

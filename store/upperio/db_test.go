package upperio_test

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/gourd/kit/store"

	"upper.io/db"
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

func testUpperDb(fn string) (adapter string, connURL db.ConnectionURL) {
	// define a database source
	adapter = sqlite.Adapter
	connURL = sqlite.ConnectionURL{
		Database: fn,
	}
	return
}

func testUpperDbData(source store.Source) (err error) {

	// dummy request
	conn, err := source.Open()
	if err != nil {
		err = fmt.Errorf("unable to connect to source: %#v",
			err.Error())
		return
	}
	defer conn.Close()

	sess := conn.Raw().(db.Database)

	drv := sess.Driver().(*sql.DB)
	_, err = drv.Exec(`
		CREATE TABLE dummy_data (
			HelloWorld text,
			FooBar text,
			Data text
		)
	`)
	if err != nil {
		err = fmt.Errorf("unable to create dummy_data: %#v",
			err.Error())
		return
	}

	coll, err := sess.Collection("dummy_data")
	if err != nil {
		err = fmt.Errorf("unable to connect to dummy_data: %#v",
			err.Error())
		return
	}

	_, err = coll.Append(&testData{
		HelloWorld: "foo bar",
		FooBar:     "foo bar",
		Data:       "something",
	})
	if err != nil {
		err = fmt.Errorf("unable to append dummy_data (foo bar): %#v",
			err.Error())
		return
	}

	_, err = coll.Append(&testData{
		HelloWorld: "foo bar 2",
		FooBar:     "foo bar 2",
		Data:       "something 2",
	})
	if err != nil {
		err = fmt.Errorf("unable to append dummy_data (foo bar 2): %#v",
			err.Error())
		return
	}

	_, err = coll.Append(&testData{
		HelloWorld: "foo bar 3",
		FooBar:     "foo bar 3",
		Data:       "something 3",
	})
	if err != nil {
		err = fmt.Errorf("unable to append dummy_data (foo bar 3): %#v",
			err.Error())
		return
	}

	return
}

// for dummy database content
type testData struct {
	HelloWorld string
	FooBar     string
	Data       string
}

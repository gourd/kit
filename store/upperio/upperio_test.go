package upperio_test

import (
	"os"
	"testing"

	"github.com/gourd/kit/store/upperio"
)

func TestSource(t *testing.T) {

	// testing database file
	fn := `./test.tmp`

	// test creating the new database
	// dummy request
	source := upperio.Source(testUpperDb(fn))
	conn, err := source()
	if err != nil {
		t.Fatalf(err.Error())
	}
	conn.Close()

	// clean up the temp database
	err = os.Remove(fn)
	if err != nil {
		t.Error(err.Error())
	}
}

package upperio_test

import (
	"github.com/gourd/kit/store/upperio"

	"net/http"
	"os"
	"testing"
)

func TestUpper(t *testing.T) {

	// testing database file
	fn := `./test.tmp`

	srcName := testUpperDb(fn)

	// dummy request
	r := &http.Request{}

	// test creating the new database
	_, err := upperio.Open(r, srcName)
	if err != nil {
		t.Error(err.Error())
	}
	upperio.Close(r, srcName)

	// clean up the temp database
	err = os.Remove(fn)
	if err != nil {
		t.Error(err.Error())
	}
}

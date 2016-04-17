package oauth2_test

import (
	"io/ioutil"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/gourd/kit/oauth2"
)

func TestLogger(t *testing.T) {
	oauth2.SetLogger(log.NewLogfmtLogger(ioutil.Discard))
}

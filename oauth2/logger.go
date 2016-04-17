package oauth2

import (
	"io/ioutil"
	"os"

	"github.com/go-kit/kit/log"
)

var msg, errMsg log.Logger

func init() {
	SetLogger(log.NewLogfmtLogger(ioutil.Discard))
	SetErrorLogger(log.NewLogfmtLogger(os.Stderr))
}

// SetLogger setup the default logger for all oauth2 operations
func SetLogger(v log.Logger) {
	msg = log.NewContext(v).WithPrefix("package", "oauth2")
}

// SetErrorLogger setup the error logger for all oauth2 operations
func SetErrorLogger(v log.Logger) {
	errMsg = log.NewContext(v).WithPrefix("package", "oauth2")
}

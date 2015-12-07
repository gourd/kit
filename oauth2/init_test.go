package oauth2_test

import (
	"bytes"
	"log"
	"os"
	"os/exec"

	"github.com/gourd/kit/store"
	"github.com/gourd/kit/store/upperio"
	"upper.io/db/sqlite"
)

const dbpath = "_test/sqlite3.db"

func init() {

	// find sqlite3 command
	sqlite3, err := exec.LookPath("sqlite3")
	if err != nil {
		log.Fatalf("error finding sqlite3 in system: %#v", err.Error())
	}
	log.Printf("sqlite3 path: %#v", sqlite3)

	// open the schema file
	file, err := os.Open("_test/schema.sqlite3.sql")
	if err != nil {
		log.Fatalf("error opening test schema: %#v", err.Error())
	}

	// initialize test database with sql file
	cmd := exec.Command(sqlite3, dbpath)
	var outstd bytes.Buffer
	var outerr bytes.Buffer

	cmd.Stdin = file
	cmd.Stdout = &outstd
	cmd.Stderr = &outerr

	if err := cmd.Run(); err != nil {
		log.Printf("output: %#v", outstd.String())
		log.Printf("error:  %#v", outerr.String())

		log.Fatalf("Failed to run sqlite command")
	}
}

func defaultTestSrc() store.Source {
	return upperio.Source(
		sqlite.Adapter, sqlite.ConnectionURL{
			Database: dbpath,
		})
}

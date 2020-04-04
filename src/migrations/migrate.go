package migrations

import (
	"database/sql"
	"github.com/DaggerJackfast/gst/src/domains"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
)

func Migrate(db *sql.DB) {
	sqlFile := path.Join(domains.RootPath, "src/migrations/start.sql")
	absPath, err := filepath.Abs(sqlFile)
	if err != nil {
		log.Fatal(err)
	}
	c, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.Fatal(err)
	}
	sqlString := string(c)
	_, err = db.Exec(sqlString)
	if err != nil {
		log.Fatal(err)
	}
}

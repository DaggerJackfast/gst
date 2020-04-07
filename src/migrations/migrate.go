package migrations

import (
	"database/sql"
	"github.com/DaggerJackfast/gst/src/domains"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"

	_ "github.com/lib/pq"
)

func Migrate(db *sql.DB) {
	sqlFile := GetSqlFile()
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

func GetSqlFile() string {
	sqlFile := path.Join(domains.RootPath, "src/migrations/start.sql")
	return sqlFile
}

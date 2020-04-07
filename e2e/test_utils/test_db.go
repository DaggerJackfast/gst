package test_utils

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"log"
)

type TestDatabase interface {
	Start()
	Stop()
	OpenDB() *sql.DB
	GetConnString() string
}
type testDatabase struct {
	dbName     string
	dbPassword string
	dbUser     string
	pool       *dockertest.Pool
	resource   *dockertest.Resource
}

func NewTestDb(dbName, dbUser, dbPassword string) TestDatabase {
	return &testDatabase{
		dbName:     dbName,
		dbPassword: dbPassword,
		dbUser:     dbUser,
	}
}

func (testDb *testDatabase) Start() {
	var db *sql.DB
	var err error
	testDb.pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to test docker: %s", err)
	}
	credentials := []string{
		fmt.Sprintf("POSTGRES_USER=%s", testDb.dbUser),
		fmt.Sprintf("POSTGRES_PASSWORD=%s", testDb.dbPassword),
		fmt.Sprintf("POSTGRES_DB=%s", testDb.dbName),
	}
	testDb.resource, err = testDb.pool.Run("postgres", "12.2", credentials)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	connString := testDb.GetConnString()
	err = testDb.pool.Retry(func() error {
		var err error

		db, err = sql.Open("postgres", connString)
		if err != nil {
			return err
		}
		return db.Ping()
	})
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not test connect to database in the docker: %s", err)
	}
	log.Print("Database in the docker has been started.")
}

func (testDb *testDatabase) Stop() {
	resource := testDb.resource
	err := testDb.pool.Purge(resource)
	if err != nil {
		log.Fatalf("Could not purge docker: %s", err)
	}
}

func (testDb *testDatabase) OpenDB() *sql.DB {
	connString := testDb.GetConnString()
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Could not test connect to docker: %s", err)
	}
	return db
}

func (testDb *testDatabase) GetConnString() string {
	if testDb.resource == nil {
		log.Fatal("Docker resource is nil. Please start database")
	}
	connString := fmt.Sprintf(
		"user=%s password=%s dbname=%s port=%s sslmode=disable",
		testDb.dbUser,
		testDb.dbPassword,
		testDb.dbName,
		testDb.resource.GetPort("5432/tcp"),
	)
	return connString
}

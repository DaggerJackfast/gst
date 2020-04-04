package spec_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	gstApp "github.com/DaggerJackfast/gst/src/app"
	"github.com/DaggerJackfast/gst/src/migrations"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"log"
	"net/http"
	"net/http/httptest"
)

func initTestDb() (*sql.DB, *dockertest.Pool, *dockertest.Resource) {
	var db *sql.DB
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to test docker: %s", err)
	}
	dbName := "test_gst"
	dbPassword := "root"
	dbUser := "root"
	credentials := []string{
		fmt.Sprintf("POSTGRES_USER=%s", dbUser),
		fmt.Sprintf("POSTGRES_PASSWORD=%s", dbPassword),
		fmt.Sprintf("POSTGRES_DB=%s", dbName),
	}
	resource, err := pool.Run("postgres", "9.6", credentials)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	err = pool.Retry(func() error {
		var err error
		//connString := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbUser, dbPassword, resource.GetPort("5432/tcp"), dbName)
		connString := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, resource.GetPort("5432/tcp"))
		db, err = sql.Open("postgres", connString)
		if err != nil {
			return err
		}
		return db.Ping()
	})
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return db, pool, resource
}



var _ = Describe("Auth", func() {
	var (
		pool     *dockertest.Pool
		resource *dockertest.Resource
		server   *httptest.Server
		client   *http.Client
		app      gstApp.Application
	)
	BeforeEach(func() {
		app = gstApp.Application{}
		app.InitLog("../../log/tests.log")
		app.Db, pool, resource = initTestDb()
		migrations.Migrate(app.Db)
		app.InitRoutes()
		server = httptest.NewServer(app.Router)
		client = &http.Client{}
	})
	AfterEach(func() {
		app.Db.Close()
		err := pool.Purge(resource)
		if err != nil {
			log.Fatalf("Could not purge docker: %s", err)
		}
		server.Close()
	})
	It("User can register", func() {
		user := map[string]string{
			"username": "test_user",
			"email":    "test@test.com",
			"password": "qweqweqwe",
		}
		data, err := json.Marshal(user)
		if err != nil {
			log.Printf("json marshal error: %s", err)
		}
		url := fmt.Sprintf("%s/auth/register", server.URL)
		response, err := client.Post(url, "application/json", bytes.NewBuffer(data))
		Expect(err).ToNot(HaveOccurred())
		Expect(response.StatusCode).To(Equal(http.StatusOK))
	})

})

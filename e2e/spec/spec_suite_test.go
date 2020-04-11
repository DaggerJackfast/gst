package spec_test

import (
	"github.com/DaggerJackfast/gst/e2e/test_utils"
	gstApp "github.com/DaggerJackfast/gst/src/app"
	"github.com/DaggerJackfast/gst/src/domains"
	"github.com/DaggerJackfast/gst/src/migrations"
	"github.com/go-testfixtures/testfixtures/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
)

func TestSpec(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GST A Suite")
}

var (
	TestDb  test_utils.TestDatabase
	Server  *httptest.Server
	Client  *http.Client
	App     gstApp.Application
	Cleaner dbcleaner.DbCleaner
	Loader  *testfixtures.Loader
)

var _ = BeforeSuite(func() {
	dbName := "test_gst"
	dbPassword := "root"
	dbUser := "root"
	App = gstApp.Application{}
	logPath := path.Join(domains.RootPath, "log/tests.log")
	App.InitLog(logPath)
	TestDb = test_utils.NewTestDb(dbName, dbUser, dbPassword)
	TestDb.Start()
	App.Db = TestDb.OpenDB()
	migrations.Migrate(App.Db)
	App.InitRoutes()
	Cleaner = dbcleaner.New()
	connString := TestDb.GetConnString()
	dbEngine := engine.NewPostgresEngine(connString)
	Cleaner.SetEngine(dbEngine)
	builder := test_utils.NewFixtureLoaderBuilder()
	Loader = builder.Build(App.Db, path.Join(domains.RootPath, "e2e/spec/fixtures"))
	Server = httptest.NewServer(App.Router)
	Client = &http.Client{}
})

var _ = AfterSuite(func() {
	Server.Close()
	App.Db.Close()
	TestDb.Stop()
})

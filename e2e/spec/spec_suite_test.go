package spec_test

import (
	"github.com/DaggerJackfast/gst/e2e/test_utils"
	gstApp "github.com/DaggerJackfast/gst/src/app"
	"github.com/DaggerJackfast/gst/src/domains"
	"github.com/DaggerJackfast/gst/src/migrations"
	"github.com/go-testfixtures/testfixtures/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
	"log"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
	"text/template"
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

	files := []string{
		path.Join(domains.RootPath, "e2e/spec/fixtures/users.yaml"),
		path.Join(domains.RootPath, "e2e/spec/fixtures/sessions.yaml"),
		path.Join(domains.RootPath, "e2e/spec/fixtures/user_profile_tokens.yaml"),
	}
	passwordValue := "qweqweqwe"
	hash, _ := bcrypt.GenerateFromPassword([]byte(passwordValue), bcrypt.MinCost)
	passwordHash := string(hash)
	var passwords []string

	for i := 0; i < 3; i++ {
		passwords = append(passwords, passwordHash)
	}
	var err error
	fixtureFunctions := template.FuncMap{
		"GenPasswordHash": func(password string) string {
			passwordBytes := []byte(password)
			hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
			if err != nil {
				log.Fatalf("Cannot generate password hash: %s", hash)
			}
			hashString := string(hash)
			return hashString
		},
	}
	Loader, err = testfixtures.New(
		testfixtures.Database(App.Db),
		testfixtures.Dialect("postgresql"),
		testfixtures.Template(),
		testfixtures.TemplateFuncs(fixtureFunctions),
		//testfixtures.TemplateData(map[string]interface{}{
		//	"Passwords": passwords,
		//}),
		testfixtures.Files(files...),
	)
	if err != nil {
		log.Fatalf("Failed in create fixtures: %s", err)
	}

	Server = httptest.NewServer(App.Router)
	Client = &http.Client{}
})

var _ = AfterSuite(func() {
	Server.Close()
	App.Db.Close()
	TestDb.Stop()
})

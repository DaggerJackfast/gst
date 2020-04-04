package spec
//
//import (
//	"database/sql"
//	gstApp "github.com/DaggerJackfast/gst/src/app"
//	"github.com/gavv/httpexpect"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"net/http/httptest"
//	"os"
//	"path/filepath"
//	"testing"
//)
//
//func TestMain(m *testing.M){
//	app := setUp()
//	retCode := m.Run()
//	tearDown(app)
//	os.Exit(retCode)
//}
//
//
//
//func setUp() *gstApp.Application {
//	app := gstApp.Application{}
//	app.InitLog("../log/tests.log")
//	app.Db = initTestDb()
//	app.InitRoutes()
//	return &app
//}
//
//func tearDown(app *gstApp.Application) {
//	app.Db.Close()
//}
//
//
//func initTestDb() *sql.DB {
//	db, err := sql.Open("ramsql", "test_gst")
//	if err != nil {
//		log.Fatal(err.Error())
//		return nil
//	}
//	migrate(db)
//	return db
//}
//
//func migrate(db *sql.DB) {
//	sqlFile := "./migrations/start.sql"
//	absPath, err := filepath.Abs(sqlFile)
//	if err != nil{
//		log.Fatal(err)
//	}
//	c, err := ioutil.ReadFile(absPath)
//	if err != nil {
//		log.Fatal(err)
//	}
//	sqlString := string(c)
//	_, err = db.Exec(sqlString)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//func TestApp(t *testing.T) {
//	server := httptest.NewServer(app.Router)
//	defer server.Close()
//	e := httpexpect.New(t, server.URL)
//	user := map[string]interface{}{
//		"username": "TestUser",
//		"email": "test@test.com",
//		"password": "qweqweqwe",
//	}
//	e.POST("/auth/register").WithJSON(user).Expect().Status(http.StatusOK)
//	tearDown(app)
//}

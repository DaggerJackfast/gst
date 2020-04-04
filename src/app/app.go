package app

import (
	"database/sql"
	"fmt"
	"github.com/DaggerJackfast/gst/src/controllers"
	"github.com/DaggerJackfast/gst/src/domains"
	"github.com/DaggerJackfast/gst/src/middlewares"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type ApplicationInterface interface {
	Initialize(user, password, dbname, logPath string)
	InitLog(logPath string)
	InitDb(user string, password string, dbname string)
	InitRoutes()
	Run(addr string)
}

type Application struct {
	Router *mux.Router
	Db     *sql.DB
	Logger *log.Logger
}

func (app *Application) Initialize(user, password, dbname, logPath string) {
	fmt.Println("The application is initializing")
	fmt.Println("The application path: ", domains.RootPath)
	// Init log
	app.InitLog(logPath)

	// Init db
	app.InitDb(user, password, dbname)
	// Init router
	app.InitRoutes()
}

func (app *Application) InitRoutes() {
	// Init controllers
	authController := controllers.NewAuthController(*app.Db, app.Logger)

	router := mux.NewRouter()
	app.Router = router
	router.HandleFunc("/auth/register", middlewares.SetMiddlewareJSON(authController.Register)).Methods("POST")
	router.HandleFunc("/auth/login", middlewares.SetMiddlewareJSON(authController.Login)).Methods("POST")
	router.HandleFunc("/auth/forgot-password", middlewares.SetMiddlewareJSON(authController.ForgotPassword)).Methods("POST")
	router.HandleFunc("/auth/change-password",
		middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(authController.ChangePassword))).Methods("POST")
	router.HandleFunc("/auth/reset-password",
		middlewares.SetMiddlewareJSON(authController.ResetPassword)).Methods("POST")
	router.HandleFunc("/auth/refresh-token",
		middlewares.SetMiddlewareJSON(authController.RefreshToken)).Methods("POST")
}

func (app *Application) InitDb(user string, password string, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		app.Logger.Fatal(err)
	}
	app.Db = db
}

func (app *Application) InitLog(logPath string) {
	absPath, err := filepath.Abs(logPath)
	if err != nil {
		log.Fatal(err)
	}
	absDirPath := filepath.Dir(absPath)
	if _, err := os.Stat(absDirPath); os.IsNotExist(err) {
		err = os.MkdirAll(absDirPath, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
	file, err := os.OpenFile(absPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	logger := log.New(file, "Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.Logger = logger
}

func (app *Application) Run(addr string) {
	fmt.Printf("The server is started at %s \n", addr)
	defer app.Db.Close()
	err := http.ListenAndServe(addr, app.Router)
	if err != nil {
		app.Logger.Fatalf("The server has stopped with the error: %s", err)
	}
}

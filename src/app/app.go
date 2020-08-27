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
	checkVariables()
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
	//check environment variables
	app.checkVariables()
	// Init db
	app.InitDb(user, password, dbname)
	// Init router
	app.InitRoutes()
}

func (app *Application) checkVariables(){
	var exists bool
	checkedEnvs := []string{"API_SECRET", "RUN_MODE"}
	for _, env := range checkedEnvs {
		_, exists = os.LookupEnv(env)
		if !exists {
			app.Logger.Fatalf("Enviroment variable '%s' is not set", env)
		}
	}
}

func (app *Application) InitRoutes() {
	// Init controllers
	authController := controllers.NewAuthController(app.Db, app.Logger)

	router := mux.NewRouter()
	app.Router = router
	router.Use(middlewares.JsonMiddleware)
	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", authController.Register).Methods("POST")
	auth.HandleFunc("/login", authController.Login).Methods("POST")
	auth.HandleFunc("/forgot-password", authController.ForgotPassword).Methods("POST")
	auth.HandleFunc("/reset-password", authController.ResetPassword).Methods("POST")
	auth.HandleFunc("/refresh-token", authController.RefreshToken).Methods("POST")

	profile := router.PathPrefix("/profile").Subrouter()
	profile.Use(middlewares.AuthenticationMiddleware)
	profile.HandleFunc("/change-password", authController.ChangePassword).Methods("POST")

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

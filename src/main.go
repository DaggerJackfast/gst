package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	dotenv "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Application interface {
	Initialize(user, password, dbname, logPath string)
	Run(addr string)
}

type App struct {
	Router 	*mux.Router
	Db *sql.DB
	Logger *log.Logger
}


func (app *App) Initialize(user, password, dbname, logPath string){
	fmt.Println("the application is initializing")
	// Init log
	absPath, err := filepath.Abs(logPath)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(absPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	logger := log.New(file, "Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.Logger = logger

	// Init db
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		app.Logger.Fatal(err)
	}
	defer db.Close()
	app.Db = db


	// Init controllers
	authController := NewAuthController(*app.Db, app.Logger)
	// Init router
	router := mux.NewRouter()
	app.Router = router
	router.HandleFunc("/auth/register", SetMiddlewareJSON(authController.Register)).Methods("POST")
	router.HandleFunc("/auth/login", SetMiddlewareJSON(authController.Login)).Methods("POST")
	router.HandleFunc("/auth/forgot-password", SetMiddlewareJSON(authController.ForgotPassword)).Methods("POST")
	router.HandleFunc("/auth/change-password",
		SetMiddlewareJSON(SetMiddlewareAuthentication(authController.ChangePassword))).Methods("POST")
	router.HandleFunc("/auth/reset-password",
		SetMiddlewareJSON(authController.ResetPassword)).Methods("POST")
	router.HandleFunc("/auth/refresh-token",
		SetMiddlewareJSON(authController.RefreshToken)).Methods("POST")
	fmt.Printf("The server is started at %s \n", "http://0.0.0.0:5050")


}

func (app *App) Run(addr string) {
	app.Logger.Fatal(http.ListenAndServe(addr, app.Router))
}

func run() {
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	logPath := os.Getenv("LOG_PATH")
	addr := os.Getenv("RUN_ADDR")
	app := App{}
	app.Initialize(user, password, dbname, logPath)
	app.Run(addr)
}

func init() {
	if err := dotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	run()
}

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

func test() {
	fmt.Println("The application is started...")
	DbConnectionString := "user=root password=root dbname=gst sslmode=disable"
	db, err := sql.Open("postgres", DbConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := NewUserRepository(*db)
	token_repo := NewUserProfileTokenRepository(*db)
	users, err := repo.All()
	if err != nil {
		panic(err)
	}
	for _, u := range users {
		fmt.Println(u.Id, u.Email, u.Username, u.Password)
	}
	n, err := repo.Find(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(n.Id, n.Email, n.Username, n.Password)
	m, err := repo.FindByEmail("ab@mail.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(m.Id, m.Email, m.Username, m.Password)
	m, err = repo.FindByUsername("ab")
	if err != nil {
		panic(err)
	}
	fmt.Println(m.Id, m.Email, m.Username, m.Password)
	ur := User{Email: "dd@mail.com", Username: "dd", Password: "ddddd"}
	user, err := repo.Store(&ur)
	if err != nil {
		panic(err)
	}
	fmt.Println(user.Id, user.Email, user.Username, user.Password)
	user.Email = "dd.new@mail.com"
	user.Username = "dd.new"
	err = repo.Update(user)
	if err != nil {
		panic(err)
	}
	fmt.Println(user.Id, user.Email, user.Username, user.Password)
	err = repo.Delete(user.Id)
	if err != nil {
		panic(err)
	}
	us := NewAuthService(repo, token_repo)
	usu := User{Email: "ee@mail.com", Username: "ee", Password: "qweqwe"}
	err = us.Register(&usu)
	if err != nil {
		panic(err)
	}
	fmt.Println(usu.Id, usu.Email, usu.Username, usu.Password)
	password := "qweqwe"
	usu.Password = password
	status := us.IsValidPassword(&usu, password)
	fmt.Println(status)
	err = repo.Delete(user.Id)
	if err != nil {
		panic(err)
	}
}

func run() {
	fmt.Println("The server is initializing...")
	db, err := sql.Open("postgres", "user=root password=root dbname=gst sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	absPath, err := filepath.Abs("../log")
	if err != nil {
		panic(err)
	}
	generalLogFile, err := os.OpenFile(absPath+"/general.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer generalLogFile.Close()

	var logger = log.New(generalLogFile, "General Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	userController := NewAuthController(*db, logger)
	router := mux.NewRouter()
	router.HandleFunc("/auth/register", SetMiddlewareJSON(userController.Register)).Methods("POST")
	router.HandleFunc("/auth/login", SetMiddlewareJSON(userController.Login)).Methods("POST")
	router.HandleFunc("/auth/forgot-password", SetMiddlewareJSON(userController.ForgotPassword)).Methods("POST")
	router.HandleFunc("/auth/change-password",
		SetMiddlewareJSON(SetMiddlewareAuthentication(userController.ChangePassword))).Methods("POST")
	router.HandleFunc("/auth/reset-password",
		SetMiddlewareJSON(userController.ResetPassword)).Methods("POST")
	router.HandleFunc("/auth/refresh-token",
		SetMiddlewareJSON(userController.RefreshToken)).Methods("POST")
	fmt.Printf("The server is started at %s \n", "http://0.0.0.0:5050")
	logger.Fatal(http.ListenAndServe(":5050", router))
}

func init() {
	if err := dotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	run()
}

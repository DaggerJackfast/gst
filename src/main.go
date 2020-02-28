package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("The application is started...")
	DbConnectionString := "user=root password=root dbname=gst sslmode=disable"
	db, err := sql.Open("postgres", DbConnectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := NewUserRepository(*db)
	users, err := repo.All()
	if err != nil {
		panic(err)
	}
	for _, u := range users{
		fmt.Println(u.Id, u.Email, u.Username, u.Password)
	}
}

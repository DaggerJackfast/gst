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
	for _, u := range users {
		fmt.Println(u.Id, u.Email, u.Username, u.Password)
	}
	n, err := repo.Find(1)
	if err != nil{
		panic(err)
	}
	fmt.Println(n.Id, n.Email, n.Username, n.Password)
	m, err := repo.FindByEmail("ab@mail.com")
	if err != nil{
		panic(err)
	}
	fmt.Println(m.Id, m.Email, m.Username, m.Password)
	m, err = repo.FindByUsername("ab")
	if err != nil{
		panic(err)
	}
	fmt.Println(m.Id, m.Email, m.Username, m.Password)
	ur := User{Email:"dd@mail.com", Username:"dd", Password:"ddddd"}
	user, err := repo.Store(&ur)
	if err != nil{
		panic(err)
	}
	fmt.Println(user.Id, user.Email, user.Username, user.Password)
	user.Email = "dd.new@mail.com"
	user.Username = "dd.new"
	err = repo.Update(user)
	if err != nil{
		panic(err)
	}
	fmt.Println(user.Id, user.Email, user.Username, user.Password)
	err = repo.Delete(user.Id)
	if err != nil{
		panic(err)
	}
}

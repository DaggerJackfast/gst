package main

import (
	"database/sql"
	"fmt"
)

type UserRepository interface {
	Find(id int) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
	All() ([]*User, error)
	Update(user *User) error
	Store(user *User) error
	Add(user *User) error
}

type userRepository struct {
	db sql.DB
}

func NewUserRepository(db sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) All() ([]User, error) {
	rows, err := repo.db.Query("select * from users order by email;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.Id, &u.Email, &u.Username, &u.Password)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, u)
	}
	return users, nil
}

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
	Store(user *User) (*User, error)
	Delete(id int64) error
}

type userRepository struct {
	db sql.DB
}

func NewUserRepository(db sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) Find(id int) (*User, error) {
	row := repo.db.QueryRow("select * from users where id=$1", id)
	var user User
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) FindByEmail(email string) (*User, error) {
	row := repo.db.QueryRow("select * from users where email=$1", email)
	var user User
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) FindByUsername(username string) (*User, error) {
	row := repo.db.QueryRow("select * from users where username=$1", username)
	var user User
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) All() ([]*User, error) {
	rows, err := repo.db.Query("select * from users order by email;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*User
	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.Id, &u.Email, &u.Username, &u.Password)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, &u)
	}
	return users, nil
}

func (repo *userRepository) Update(user *User) error {
	_, err := repo.db.Exec("update users set email=$2, username=$3, password=$4 where id=$1",
		user.Id, user.Email, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) Store(user *User) (*User, error) {
	err := repo.db.QueryRow("insert into users (email, username, password) values($1, $2, $3) returning  id",
		user.Email, user.Username, user.Password).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) Delete(id int64) error {
	_, err := repo.db.Exec("delete from users where id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

package repositories

import (
	"database/sql"
	"github.com/DaggerJackfast/gst/src/domains"
)

type UserRepository interface {
	Find(id uint64) (*domains.User, error)
	FindByEmail(email string) (*domains.User, error)
	FindByUsername(username string) (*domains.User, error)
	All() ([]*domains.User, error)
	Update(user *domains.User) error
	Store(user *domains.User) (*domains.User, error)
	Delete(id uint64) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) Find(id uint64) (*domains.User, error) {
	row := repo.db.QueryRow("select * from users where id=$1", id)
	var user domains.User
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) FindByEmail(email string) (*domains.User, error) {
	row := repo.db.QueryRow("select * from users where email=$1", email)
	var user domains.User
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) FindByUsername(username string) (*domains.User, error) {
	row := repo.db.QueryRow("select * from users where username=$1", username)
	var user domains.User
	err := row.Scan(&user.Id, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) All() ([]*domains.User, error) {
	rows, err := repo.db.Query("select * from users order by email;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*domains.User
	for rows.Next() {
		u := domains.User{}
		err := rows.Scan(&u.Id, &u.Email, &u.Username, &u.Password)
		if err != nil {
			continue
		}
		users = append(users, &u)
	}
	return users, nil
}

func (repo *userRepository) Update(user *domains.User) error {
	_, err := repo.db.Exec("update users set email=$2, username=$3, password=$4 where id=$1",
		user.Id, user.Email, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) Store(user *domains.User) (*domains.User, error) {
	err := repo.db.QueryRow("insert into users (email, username, password) values($1, $2, $3) returning  id",
		user.Email, user.Username, user.Password).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) Delete(id uint64) error {
	_, err := repo.db.Exec("delete from users where id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

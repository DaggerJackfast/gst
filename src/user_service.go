package main

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserService interface {
	Register(user *User) (*User, error)
	//ForgotPassword(user *User) error
	//ChangePassword(user *User, password string) error
	//Validate(user *User)
	//Auth(user *User, password string) error
	//IsValid(user *User) bool
	//GetRepo() UserRepository
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (service *userService) Register(user *User) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hash)
	user, err = service.repo.Store(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePlainPwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

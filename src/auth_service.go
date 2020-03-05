package main

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(user *User) (*User, error)
	//ForgotPassword(user *User) error
	//Validate(user *User)
	Login(user *User) error
	ChangePassword(user *User, password string) error
	IsValidPassword(user *User, password string) bool
	GetRepo() UserRepository
}

type authService struct {
	repo UserRepository
}

func NewAuthService(repo UserRepository) AuthService {
	return &authService{repo: repo}
}

func (service *authService) Register(user *User) (*User, error) {
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

func (service *authService) Login(user *User) error {
	password := user.Password
	user, err := service.repo.FindByEmail(user.Email)
	if err != nil {
		return err
	}
	valid := service.IsValidPassword(user, password)
	if !valid {
		return errors.New("Password is incorrect")
	}
	return nil
}

func (service *authService) ChangePassword(user *User, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	err = service.repo.Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (service *authService) IsValidPassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(err.Error())
	}
	return err == nil
}

func (service *authService) GetRepo() UserRepository {
	return service.repo
}

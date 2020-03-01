package main

import (
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *User) (*User, error)
	//ForgotPassword(user *User) error
	ChangePassword(user *User, password string) error
	//Validate(user *User)
	//Login(user *User, password string) error
	IsValidPassword(user *User) (bool, error)
	GetRepo() UserRepository
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

func (service *userService) ChangePassword(user *User, password string) error {
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

func (service *userService) IsValidPassword(user *User) (bool, error) {
	us, err := service.repo.FindByEmail(user.Email)
	if err != nil{
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(us.Password) ,[]byte(user.Password))
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (service *userService) GetRepo() UserRepository{
	return service.repo
}

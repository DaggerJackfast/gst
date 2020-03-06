package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type UserController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type userController struct {
	db     sql.DB
	logger *log.Logger
}

func NewUserController(db sql.DB, logger *log.Logger) UserController {
	return &userController{
		db:     db,
		logger: logger,
	}
}

func (controller userController) Register(w http.ResponseWriter, r *http.Request) {
	service := NewAuthService(NewUserRepository(controller.db))
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = service.Register(&user)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	JSON(w, http.StatusOK, user)
}

func (controller userController) Login(w http.ResponseWriter, r *http.Request) {
	service := NewAuthService(NewUserRepository(controller.db))
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = service.Login(&user)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusForbidden, err)
		return
	}
	token, err := CreateToken(user.Id)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	data := AuthUserToken{User: user, Token: token}
	JSON(w, http.StatusOK, data)
}

package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
	ForgotPassword(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

type authController struct {
	db     sql.DB
	logger *log.Logger
}

func NewAuthController(db sql.DB, logger *log.Logger) AuthController {
	return &authController{
		db:     db,
		logger: logger,
	}
}

func (controller authController) Register(w http.ResponseWriter, r *http.Request) {
	service := NewAuthService(NewUserRepository(controller.db), NewUserProfileTokenRepository(controller.db))
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

func (controller authController) Login(w http.ResponseWriter, r *http.Request) {
	service := NewAuthService(NewUserRepository(controller.db), NewUserProfileTokenRepository(controller.db))
	epf := EmailPasswordFingerprint{}
	err := json.NewDecoder(r.Body).Decode(&epf)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	validate := validator.New()
	err = validate.Struct(epf)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := User{
		Email: epf.Email,
		Password: epf.Password,
	}
	err = service.Login(&user)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusForbidden, err)
		return
	}
	tokenPair, err := CreateTokenPair(user.Id)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	session := Session{
		User:         &user,
		RefreshToken: tokenPair["refresh_token"],
		UserAgent:    GetUserAgent(r),
		FingerPrint:  epf.FingerPrint,
		Ip:           GetIp(r),
		ExpiredIn:    ExpiredInRefreshToken,
	}
	sessionService := NewSessionService(NewSessionRepository(controller.db))
	err = sessionService.CreateSession(&session)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	data := AuthUserToken{User: user, Token: tokenPair}
	JSON(w, http.StatusOK, data)
}
func (controller authController) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	service := NewAuthService(NewUserRepository(controller.db), NewUserProfileTokenRepository(controller.db))
	email := UserEmail{}
	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// TODO: need to delete code duplications
	validate := validator.New()
	err = validate.Struct(email)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := service.ForgotPassword(email.Email)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	em := NewEmailSender(*controller.logger)
	recipients := []string{email.Email}
	err = em.Send(recipients, "root@root.root", fmt.Sprintf("token: %s", token.ProfileToken))
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	data := Response{Status: Success, Message: "Please check your email"}
	JSON(w, http.StatusOK, data)
}

func (controller authController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	service := NewAuthService(NewUserRepository(controller.db), NewUserProfileTokenRepository(controller.db))
	ept := EmailPasswordToken{}
	err := json.NewDecoder(r.Body).Decode(&ept)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	validate := validator.New()
	err = validate.Struct(ept)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = service.ResetPassword(ept.Email, ept.Password, ept.Token)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	data := Response{Status: Success, Message: "Your password successfully changed."}
	JSON(w, http.StatusOK, data)
}

func (controller authController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	service := NewAuthService(NewUserRepository(controller.db), NewUserProfileTokenRepository(controller.db))
	userId, err := ExtractTokenId(r)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnauthorized, err)
		return
	}
	currentUser, err := service.GetUser(userId)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusNotFound, err)
		return
	}
	p := Passwords{}
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	validate := validator.New()
	err = validate.Struct(p)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	if valid := service.IsValidPassword(currentUser, p.OldPassword); !valid {
		err = errors.New("Password is incorrect")
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusForbidden, err)
		return
	}
	err = service.ChangePassword(currentUser, p.NewPassword)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	JSON(w, http.StatusOK, currentUser)
}

func (controller authController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	service := NewAuthService(NewUserRepository(controller.db), NewUserProfileTokenRepository(controller.db))
	sessionService := NewSessionService(NewSessionRepository(controller.db))

	rToken := RefreshToken{}
	err := json.NewDecoder(r.Body).Decode(&rToken)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	validate := validator.New()
	err = validate.Struct(rToken)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userId, err := ExtractId(rToken.RefreshToken)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusForbidden, err)
		return
	}
	user, err := service.GetUser(userId)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusForbidden, err)
		return
	}
	err = service.Login(user)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusForbidden, err)
		return
	}
	tokenPair, err := CreateTokenPair(user.Id)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	session, err := sessionService.GetSession(user, rToken.RefreshToken)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusForbidden, err)
		return
	}
	session.RefreshToken = tokenPair["refresh_token"]
	err = sessionService.UpdateSession(session)
	if err != nil {
		controller.logger.Println(err.Error())
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	data := AuthUserToken{User: *user, Token: tokenPair}
	JSON(w, http.StatusOK, data)
}

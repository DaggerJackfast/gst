package services

import (
	"errors"
	"fmt"
	"github.com/DaggerJackfast/gst/src/domains"
	"github.com/DaggerJackfast/gst/src/repositories"
	gstToken "github.com/DaggerJackfast/gst/src/token"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService interface {
	Register(user *domains.User) error
	GetUser(userId uint64) (*domains.User, error)
	ForgotPassword(email string) (*domains.UserProfileToken, error)
	ResetPassword(email, password, token string) error
	//Validate(user *User)
	ValidateToken(user *domains.User, tokenValue, tokenType string) error
	Login(user *domains.User) error
	ChangePassword(user *domains.User, password string) error
	IsValidPassword(user *domains.User, password string) bool
	GetRepo() repositories.UserRepository
}

type authService struct {
	userRepo  repositories.UserRepository
	tokenRepo repositories.UserProfileTokenRepository
}

func NewAuthService(userRepo repositories.UserRepository, tokenRepo repositories.UserProfileTokenRepository) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (service *authService) Register(user *domains.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	us, err := service.userRepo.Store(user)
	if err != nil {
		return err
	}
	user.Modify(*us)
	return nil
}

func (service *authService) Login(user *domains.User) error {
	password := user.Password
	us, err := service.userRepo.FindByEmail(user.Email)
	if err != nil {
		return err
	}
	valid := service.IsValidPassword(us, password)
	if !valid {
		return errors.New("Password is incorrect")
	}
	user.Modify(*us)
	return nil
}

func (service *authService) ChangePassword(user *domains.User, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	err = service.userRepo.Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (service *authService) IsValidPassword(user *domains.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(err.Error())
	}
	return err == nil
}

func (service *authService) ForgotPassword(email string) (*domains.UserProfileToken, error) {
	user, err := service.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	randToken, err := gstToken.GenerateToken(16)
	if err != nil {
		return nil, err
	}
	token := domains.UserProfileToken{
		User:         user,
		ProfileToken: randToken,
		TokenType:    domains.ForgotPasswordToken,
		IsActive:     true,
		ExpiredIn:    domains.ExpiredInForgotPasswordToken,
	}
	err = service.tokenRepo.Store(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (service *authService) ResetPassword(email, password string, token string) error {
	user, err := service.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	err = service.ValidateToken(user, token, domains.ForgotPasswordToken)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	err = service.userRepo.Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (service *authService) ValidateToken(user *domains.User, tokenValue, tokenType string) error {
	token, err := service.tokenRepo.FindUserTokenByStatus(user, tokenType)
	if err != nil {
		return err
	}
	if !token.CreatedAt.Add(time.Duration(token.ExpiredIn) * time.Second).Before(time.Now()) {
		return errors.New("The token time is expired.")
	}
	if token.ProfileToken != tokenValue {
		return errors.New("The token is wrong.")
	}
	token.IsActive = false
	err = service.tokenRepo.Update(token)
	if err != nil {
		return err
	}
	return nil
}

func (service *authService) GetUser(userId uint64) (*domains.User, error) {
	user, err := service.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *authService) GetRepo() repositories.UserRepository {
	return service.userRepo
}

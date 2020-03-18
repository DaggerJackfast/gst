package main

import (
	"time"
)

type User struct {
	Id       uint64 `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (user *User) Modify(us User) {
	user.Id = us.Id
	user.Email = us.Email
	user.Username = us.Username
	user.Password = us.Password
}

type UserProfileToken struct {
	Id           uint64
	User         *User
	ProfileToken string
	TokenType    string
	IsActive     bool
	ExpiredIn    uint64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Session struct {
	Id uint64
	User *User
	RefreshToken string
	UserAgent string
	FingerPrint string
	Ip string
	ExpiredIn uint64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AuthUserToken struct {
	User  User   `json:"user"`
	Token map[string]string `json:"token"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required,omitempty" structs:"required,omitempty"`
}

type UserEmail struct {
	Email string `json:"email" validate:"email,required,omitempty" structs:"required,omitempty"`
}

type EmailPasswordFingerprint struct {
	Email    string `json:"email" validate:"required,omitempty" structs:"required,omitempty"`
	Password string `json:"password" validate:"required,omitempty" structs:"required,omitempty"`
	FingerPrint    string `json:"fingerprint" validate:"required,omitempty" structs:"required,omitempty"`
}

type EmailPasswordToken struct {
	Email    string `json:"email" validate:"required,omitempty" structs:"required,omitempty"`
	Password string `json:"password" validate:"required,omitempty" structs:"required,omitempty"`
	Token    string `json:"token" validate:"required,omitempty" structs:"required,omitempty"`
}

type Passwords struct {
	OldPassword string `json:"old_password" validate:"required,omitempty" structs:"required,omitempty"`
	NewPassword string `json:"new_password" validate:"required,omitempty" structs:"required,omitempty"`
}

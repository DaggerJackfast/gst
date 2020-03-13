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

type AuthUserToken struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type UserEmail struct {
	Email string `json:"email" validate:"email,required,omitempty" structs:"required,omitempty"`
}

type EmailPasswordToken struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Passwords struct {
	OldPassword string `json:"old_password" validate:"required,omitempty" structs:"required,omitempty"`
	NewPassword string `json:"new_password" validate:"required,omitempty" structs:"required,omitempty"`
}

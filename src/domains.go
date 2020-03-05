package main

type User struct {
	Id       uint64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthUserToken struct {
	User User `json:"user"`
	Token string `json:"token"`
}

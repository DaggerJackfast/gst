package main

type User struct {
	Id       uint64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (user *User) Modify(us User){
	user.Id = us.Id
	user.Email = us.Email
	user.Username = us.Username
	user.Password = us.Password
}

type AuthUserToken struct {
	User User `json:"user"`
	Token string `json:"token"`
}

type Passwords struct{
	OldPassword string `json:"old_password" validate:"required,omitempty" structs:"required,omitempty"`
	NewPassword string `json:"new_password" validate:"required,omitempty" structs:"required,omitempty"`
}

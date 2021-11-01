package model

import "time"

const (
	UserTypeAdmin         = 1
	UserTypeExternalAdmin = 2
	UserTypeExternalMedic = 3
)

type UserType int8

type User struct {
	ID        string    `json:"id"`
	UserName  string    `json:"user_name"`
	Password  string    `json:"password"`
	Type      UserType  `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Users []User

type UserLogin struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

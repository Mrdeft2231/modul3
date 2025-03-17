package model

import "time"

type User struct {
	Id              int
	Login           string
	Email           string
	Password        string
	OldPassword     string
	Status          int
	Role            string
	Block           int
	PasswordChanged bool
	CreateUser      time.Time
}

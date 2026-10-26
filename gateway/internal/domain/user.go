package domain

import "errors"

type User struct {
	ID           string `json:"id"`
	DisplayName  string `json:"displayName"`
	Role         string `json:"role"`
	PasswordHash string `json:"-"`
}

var ErrLoginTaken = errors.New("login already exists")

var ErrUserNotFound = errors.New("user not found")

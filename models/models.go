package models

import (
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

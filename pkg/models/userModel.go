package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Username string
	Password string
}

type SignUpDto struct {
	Name     string `validate:"required"`
	Username string `validate:"required,min=4"`
	Password string `validate:"required,min=4"`
}

type SignInDto struct {
	Username string `validate:"required,min=4"`
	Password string `validate:"required,min=4"`
}

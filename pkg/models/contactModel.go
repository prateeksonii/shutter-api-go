package models

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	Label  string
	UserID uint
}

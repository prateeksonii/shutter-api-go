package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name            string
	Username        string
	Password        string
	Contacts        []Contact
	SentInvites     []Invite `gorm:"foreignKey:SenderID"`
	ReceivedInvites []Invite `gorm:"foreignKey:ReceiverID"`
}

type SignUpDto struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required,min=4"`
	Password string `json:"password" binding:"required,min=4"`
}

type SignInDto struct {
	Username string `json:"username" binding:"required,min=4"`
	Password string `json:"password" binding:"required,min=4"`
}

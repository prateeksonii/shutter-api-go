package models

import (
	"gorm.io/gorm"
)

type Invite struct {
	gorm.Model
	SenderID   uint
	ReceiverID uint
	Message    string
}

type SendInviteDto struct {
	Message string `validate:"required"`
}

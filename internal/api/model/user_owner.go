package model

import "gorm.io/gorm"

type UserOwnedEntity struct {
	gorm.Model
	UserID string `json:"user_id" db:"user_id"`
}

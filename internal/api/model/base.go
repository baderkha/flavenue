package model

import (
	"time"

	"github.com/baderkha/library/pkg/ptr"
	"github.com/gofrs/uuid"
)

// Base : attach this as a base model using uuid
type Base struct {
	ID        string     `json:"id" db:"id" gorm:"type:VARCHAR(100);primary"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	IsDeleted bool       `json:"is_deleted" db:"is_deleted" gorm:"type:TINYINT(1);index"`
}

func (b *Base) New() {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	b.ID = id.String()
	b.IsDeleted = false
	b.CreatedAt = ptr.Get(time.Now())
	b.UpdatedAt = ptr.Get(time.Now())
}

func (b Base) GetID() string {
	return b.ID
}

func (b Base) GetIDKey() string {
	return b.ID
}

// BaseOwned : use this for entities that are owned and need to have an account to own them
type BaseOwned struct {
	Base
	AccountID string `json:"account_id" db:"account_id" gorm:"type:VARCHAR(255);index"`
}

func (b BaseOwned) GetAccountID() string {
	return b.AccountID
}

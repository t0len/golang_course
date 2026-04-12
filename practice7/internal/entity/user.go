package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `json:"ID" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt time.Time      `json:"CreatedAt"`
	UpdatedAt time.Time      `json:"UpdatedAt"`
	DeletedAt gorm.DeletedAt `json:"DeletedAt" gorm:"index"`
	Username  string         `json:"Username" gorm:"unique"`
	Email     string         `json:"Email" gorm:"unique"`
	Password  string         `json:"Password"`
	Role      string         `json:"Role"` // "user" or "admin"
	Verified  bool           `json:"Verified"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

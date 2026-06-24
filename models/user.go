package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Name            string     `gorm:"type:varchar(100);not null" json:"name" binding:"required"`
	Email           string     `gorm:"type:varchar(191);uniqueIndex;not null" json:"email" binding:"required,email"`
	Password        string     `gorm:"type:varchar(255);not null" json:"-"`
	ResetToken      *string    `gorm:"type:varchar(255);index" json:"-"`
	ResetTokenExpiry *time.Time `json:"-"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Hook do GORM para garantir a geração do UUID antes de criar o registro
// Block para o caso de o banco não usar default nativo:
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

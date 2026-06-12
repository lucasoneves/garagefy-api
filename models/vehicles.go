package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Vehicle struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Brand        string         `gorm:"type:varchar(100);not null" json:"brand"`
	Model        string         `gorm:"type:varchar(100);not null" json:"model"`
	Year         int            `gorm:"not null" json:"year"`
	Plate        string         `gorm:"type:varchar(20);unique;not null" json:"plate"`
	CurrentOdo   int            `gorm:"not null" json:"current_odo"`
	Color        string         `gorm:"type:varchar(50)" json:"color"`
	LogbookLines []LogbookEntry `gorm:"foreignKey:VehicleID;constraint:OnDelete:CASCADE;" json:"logbook_lines,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate garante a geração automática do UUID antes de salvar no banco
func (v *Vehicle) BeforeCreate(tx *gorm.DB) (err error) {
	v.ID = uuid.New()
	return
}

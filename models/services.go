package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"id"`
	VehicleID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"vehicle_id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	ShopName    string         `json:"shop_name"`
	CurrentOdo  int            `gorm:"not null" json:"current_odo"`
	Cost        float64        `gorm:"type:numeric(10,2);not null" json:"cost"`
	ServiceDate time.Time      `gorm:"not null" json:"service_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Service) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}

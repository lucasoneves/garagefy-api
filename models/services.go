package models

import (
	"time"
)

type Service struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	VehicleID   string    `gorm:"type:uuid;not null;index" json:"vehicle_id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	ShopName    string    `json:"shop_name"`
	CurrentOdo  int       `gorm:"not null" json:"current_odo"`
	Cost        float64   `gorm:"type:numeric(10,2);not null" json:"cost"`
	ServiceDate time.Time `gorm:"not null" json:"service_date"`
	CreatedAt   time.Time `json:"created_at"`
}

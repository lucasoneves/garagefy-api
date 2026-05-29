package models

import (
	"time"
)

type LogbookCategory string

const (
	Observation LogbookCategory = "Observation"
	Reminder    LogbookCategory = "Reminder"
	Todo        LogbookCategory = "To-do"
)

type LogbookEntry struct {
	ID            string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	VehicleID     string          `gorm:"type:uuid;not null;index" json:"vehicle_id"`
	Category      LogbookCategory `gorm:"type:varchar(20);not null" json:"category"`
	Title         string          `gorm:"type:varchar(255);not null" json:"title"`
	Description   string          `gorm:"type:text;not null" json:"description"`
	AttachmentURL *string         `gorm:"type:varchar(512)" json:"attachment_url,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

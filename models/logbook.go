package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LogbookCategory string

const (
	Observation LogbookCategory = "Observation"
	Reminder    LogbookCategory = "Reminder"
	Todo        LogbookCategory = "To-do"
)

type LogbookEntry struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey;" json:"id"`
	VehicleID     uuid.UUID       `gorm:"type:uuid;not null;index" json:"vehicle_id"` // <-- Garanta esta linha
	Category      LogbookCategory `gorm:"type:varchar(50);not null" json:"category"`
	Title       string          `gorm:"type:varchar(255);not null" json:"title"`
	Description string          `gorm:"type:text" json:"description"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}

// BeforeCreate garante a geração automática do UUID para o logbook antes de salvar
func (l *LogbookEntry) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID = uuid.New()
	return
}

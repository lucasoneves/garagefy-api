package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FuelLog struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;" json:"id"`
	VehicleID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"vehicle_id"`
	Date        time.Time      `gorm:"not null" json:"date"`
	Odometer    int            `gorm:"not null" json:"odometer"`
	GasStation  string         `gorm:"type:varchar(100)" json:"gas_station,omitempty"`
	FuelType    string         `gorm:"type:varchar(20);not null;default:'Gasolina'" json:"fuel_type"`
	Liters      float64        `gorm:"type:numeric(6,2);not null" json:"liters"`
	PricePerLit float64        `gorm:"type:numeric(5,2);not null" json:"price_per_liter"`
	TotalCost   float64        `gorm:"type:numeric(8,2)" json:"total_cost"`
	IsFullTank  bool           `gorm:"default:false" json:"is_full_tank"`
	KmLiter     float64        `gorm:"type:numeric(5,2)" json:"km_liter,omitempty"` // Calculado automaticamente
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate garante a geração do UUID e calcula o custo total antes de salvar
func (f *FuelLog) BeforeCreate(tx *gorm.DB) (err error) {
	f.ID = uuid.New()

	// Se o custo total não for enviado, calcula automaticamente
	if f.TotalCost == 0 && f.Liters > 0 && f.PricePerLit > 0 {
		f.TotalCost = f.Liters * f.PricePerLit
	}
	return
}

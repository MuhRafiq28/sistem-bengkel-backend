package models

import (
	"time"
	"github.com/lib/pq"	
)

type Product struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	Name      string        `json:"name"`
	Brand     string        `json:"brand"`
	Price     float64       `json:"price"`
	Stock     int           `json:"stock"`
	Gram      *int          `json:"gram,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	RPM       pq.Int64Array `gorm:"type:integer[]" json:"rpm,omitempty"`
	Volume    *string       `json:"volume,omitempty"`
}


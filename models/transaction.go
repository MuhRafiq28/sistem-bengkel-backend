package models

import "time"
import "github.com/lib/pq"

type Transaction struct {
	ID        uint                `gorm:"primaryKey" json:"id"`
	Date      time.Time           `json:"date"`
	Total     float64             `json:"total"`
	Details   []TransactionDetail `json:"details"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type TransactionDetail struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	TransactionID uint          `json:"transaction_id"`
	ProductID     *uint         `json:"product_id"`
	ProductName   string        `json:"product_name"`
	Brand         string        `json:"brand"`
	Gram          *int          `json:"gram,omitempty"`
	Volume        *string       `json:"volume,omitempty"`
	RPM           pq.Int64Array `gorm:"type:integer[]" json:"rpm,omitempty"`
	Qty           int           `json:"qty"`
	Price         float64       `json:"price"`
	Subtotal      float64       `json:"subtotal"`
}


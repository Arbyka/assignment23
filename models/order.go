package models

import (
	"time"
)

type Order struct {
	OrderID   uint      `gorm:"primaryKey" json:"order_id"`
	ProductID uint      `gorm:"not null" json:"product_id"`
	Quantity  int       `json:"quantity"`
	OrderDate time.Time `json:"order_date"`
	Product   Product   `gorm:"foreignKey:ProductID"` // Relasi dengan tabel produk
}

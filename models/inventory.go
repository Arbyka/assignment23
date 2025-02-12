package models

type Inventory struct {
	InventoryID uint    `gorm:"primaryKey" json:"inventory_id"`
	ProductID   uint    `gorm:"not null" json:"product_id"`
	Quantity    int     `json:"quantity"`
	Location    string  `json:"location"`
	Product     Product `gorm:"foreignKey:ProductID"` // Relasi dengan produk
}

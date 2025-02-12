package controllers

import (
	"golangapi/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InventoryController struct {
	DB *gorm.DB
}

// Constructor untuk InventoryController
func NewInventoryController(db *gorm.DB) *InventoryController {
	return &InventoryController{DB: db}
}

// Get All Inventory
func (pc *InventoryController) GetInventory(c *gin.Context) {
	var inventory []models.Inventory
	pc.DB.Unscoped().Find(&inventory)
	c.JSON(http.StatusOK, gin.H{"inventory": inventory})
}

// Get Stock Level
func (ic *InventoryController) GetStock(c *gin.Context) {
	productID := c.Param("product_id")

	var inventory models.Inventory
	if err := ic.DB.Where("product_id = ?", productID).First(&inventory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stok tidak ditemukan untuk produk ini"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"inventory": inventory})
}

// Update Stock Level
func (ic *InventoryController) UpdateStock(c *gin.Context) {
	productID := c.Param("product_id")
	quantityChange, err := strconv.Atoi(c.Query("quantity")) // Ambil query param quantity

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jumlah stok harus berupa angka"})
		return
	}

	var inventory models.Inventory
	if err := ic.DB.Where("product_id = ?", productID).First(&inventory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stok tidak ditemukan"})
		return
	}

	// Update stok
	inventory.Quantity += quantityChange
	ic.DB.Save(&inventory)

	c.JSON(http.StatusOK, gin.H{"message": "Stok berhasil diperbarui", "inventory": inventory})
}

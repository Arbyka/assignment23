package controllers

import (
	"golangapi/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderController struct {
	DB *gorm.DB
}

// Constructor untuk OrderController
func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{DB: db}
}

// Menampilkan semua pesanan
func (oc *OrderController) GetOrders(c *gin.Context) {
	var orders []models.Order
	oc.DB.Preload("Product", "product_id = ?", "product_id").Find(&orders)
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// Mengambil detail pesanan berdasarkan ID
func (oc *OrderController) GetOrderByID(c *gin.Context) {
	orderID := c.Param("order_id")

	var order models.Order
	if err := oc.DB.Preload("Product").First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// Membuat pesanan baru
func (oc *OrderController) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set tanggal pesanan ke waktu saat ini
	order.OrderDate = time.Now()

	// Simpan pesanan ke database
	oc.DB.Create(&order)
	c.JSON(http.StatusCreated, gin.H{"message": "Pesanan berhasil dibuat", "order": order})
}

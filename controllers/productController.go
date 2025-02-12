package controllers

import (
	"golangapi/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
	DB *gorm.DB
}

// Constructor untuk ProductController
func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{DB: db}
}

// Get All Products
func (pc *ProductController) GetProducts(c *gin.Context) {
	var products []models.Product
	pc.DB.Unscoped().Find(&products)
	c.JSON(http.StatusOK, gin.H{"products": products})
}

// Get Product by ID
func (pc *ProductController) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := pc.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}

// Get Products by Category
func (pc *ProductController) GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")
	var products []models.Product
	pc.DB.Where("category = ?", category).Find(&products)
	c.JSON(http.StatusOK, gin.H{"products": products})
}

// Create New Product
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pc.DB.Create(&product)
	c.JSON(http.StatusCreated, gin.H{"product": product})
}

// Update Product
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	// Cari produk berdasarkan ID
	if err := pc.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	// Bind JSON ke struct baru agar tidak menimpa ID
	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update hanya field yang diperlukan
	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Price = updatedProduct.Price

	pc.DB.Save(&product)
	c.JSON(http.StatusOK, gin.H{"product": product})
}

// Delete Product
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := pc.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	pc.DB.Delete(&product)
	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
}

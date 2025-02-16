package controllers

import (
	"os"
	"path/filepath"
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

// Upload Product Image
func (pc *ProductController) UploadProductImage(c *gin.Context) {
	id := c.Param("id")

	// Cari produk berdasarkan ID
	var product models.Product
	if err := pc.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	// Validasi dan ambil file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan"})
		return
	}

	// Validasi format file
	allowedExtensions := map[string]bool{".png": true, ".jpg": true, ".jpeg": true}
	ext := filepath.Ext(file.Filename)
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak valid, hanya PNG, JPG, dan JPEG diperbolehkan"})
		return
	}

	// Buat folder uploads/products jika belum ada
	uploadDir := "uploads/products"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat direktori penyimpanan"})
		return
	}

	// Simpan file dengan nama ID produk
	filename := id + ext
	filePath := filepath.Join(uploadDir, filename)

	// Simpan file ke disk
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
		return
	}

	// Update path gambar di database
	product.Image = filePath
	pc.DB.Save(&product)

	c.JSON(http.StatusOK, gin.H{"message": "Gambar berhasil diunggah", "file_path": filePath})
}

// Download Product Image
func (pc *ProductController) DownloadProductImage(c *gin.Context) {
	id := c.Param("id")

	// Cari produk berdasarkan ID
	var product models.Product
	if err := pc.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	// Cek apakah produk memiliki gambar
	if product.Image == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gambar tidak tersedia"})
		return
	}

	// Kirim gambar ke client
	c.File(product.Image)
}
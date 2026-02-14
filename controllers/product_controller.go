package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"

	"github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

// Input struct untuk POST / PUT
type ProductInput struct {
	Name   string  `json:"name" binding:"required"`
	Brand  string  `json:"brand" binding:"required"`
	Price  float64 `json:"price" binding:"required"`
	Stock  int     `json:"stock" binding:"required"`
	Gram   *int    `json:"gram,omitempty"`   // optional, hanya untuk Rorer
	RPM    []int   `json:"rpm,omitempty"`    // Persentri
	Volume *string `json:"volume,omitempty"` // Oli Mesin
}

// GetProducts - ambil semua produk
func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := config.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": products})
}

// CreateProduct - tambah produk baru
func CreateProduct(c *gin.Context) {
	var input ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert []int (input.RPM) ke pq.Int64Array
	var rpm pq.Int64Array
	for _, r := range input.RPM {
		rpm = append(rpm, int64(r))
	}

	product := models.Product{
		Name:   input.Name,
		Brand:  input.Brand,
		Price:  input.Price,
		Stock:  input.Stock,
		Gram:   input.Gram,
		RPM:    rpm, // <-- pakai rpm yang sudah di-convert
		Volume: input.Volume,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": product})
}

// UpdateProduct - edit produk by ID
func UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.Name = input.Name
	product.Brand = input.Brand
	product.Price = input.Price
	product.Stock = input.Stock
	product.Gram = input.Gram

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// DeleteProduct - hapus produk by ID (hard delete)
func DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Hard delete: data benar-benar dihapus dari database
	if err := config.DB.Unscoped().Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product permanently deleted ✅",
	})
}

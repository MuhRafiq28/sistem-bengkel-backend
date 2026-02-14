package controllers

import (
	"fmt"
	"net/http"
	"time"
	"backend/config" 
	"backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*
========================================
REQUEST STRUCT
========================================
*/

type TransactionRequest struct {
	Details []TransactionDetailRequest `json:"details" binding:"required"`
}

type TransactionDetailRequest struct {
	ProductID *uint   `json:"product_id"`
	Qty       int     `json:"qty" binding:"required,min=1"`
	Price     float64 `json:"price" binding:"required"`
}

/*
========================================
CREATE TRANSACTION
========================================
*/

func CreateTransaction(c *gin.Context) {
	var req TransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("BindJSON error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	var total float64
	var details []models.TransactionDetail

	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Recovered panic:", r)
		}
	}()

	for _, item := range req.Details {
		subtotal := item.Price * float64(item.Qty)

		var detail models.TransactionDetail

		if item.ProductID != nil {

			var product models.Product

			if err := tx.First(&product, *item.ProductID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Product tidak ditemukan"})
				return
			}

			if product.Stock <= 0 {
    tx.Rollback()
    c.JSON(http.StatusBadRequest, gin.H{
        "error": "Maaf stok tidak tersedia untuk " + product.Name,
    })
    return
}

//kurangi stok produk
if product.Stock <= 0 {
    tx.Rollback()
    c.JSON(http.StatusBadRequest, gin.H{
        "error": "Maaf stok tidak tersedia untuk " + product.Name,
    })
    return
}

// ==========================
// 🔥 UPDATE STOK DI SINI
// ==========================
product.Stock -= item.Qty

if err := tx.Save(&product).Error; err != nil {
    tx.Rollback()
    c.JSON(http.StatusInternalServerError, gin.H{
        "error": "Gagal update stok",
    })
    return
}



			

			detail = models.TransactionDetail{
	ProductID:   item.ProductID,
	ProductName: product.Name,
	Brand:       product.Brand,
	Gram:        product.Gram,
	Volume:      product.Volume,
	RPM:         product.RPM,
	Qty:         item.Qty,
	Price:       item.Price,
	Subtotal:    subtotal,
}


		} else {
			detail = models.TransactionDetail{
				Qty:      item.Qty,
				Price:    item.Price,
				Subtotal: subtotal,
			}
		}

		details = append(details, detail)
		total += subtotal
	}

	transaction := models.Transaction{
		Date:    time.Now(),
		Total:   total,
		Details: details,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat transaksi"})
		return
	}

	if err := tx.Preload("Details").First(&transaction, transaction.ID).Error; err != nil {
		fmt.Println("Preload error:", err)
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":     "Transaksi berhasil",
		"transaction": transaction,
	})
}

/*
========================================
GET ALL TRANSACTIONS (HISTORY)
========================================
*/

func GetTransactionsHistory(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var transactions []models.Transaction

	if err := db.Preload("Details").
		Order("date desc").
		Find(&transactions).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil transaksi",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}

/*
========================================
DAILY REPORT
========================================
*/

func GetDailyReport(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var results []struct {
		Label string  `json:"label"`
		Total float64 `json:"total"`
	}

	db.Raw(`
		SELECT DATE(date) as label, SUM(total) as total
		FROM transactions
		GROUP BY DATE(date)
		ORDER BY DATE(date)
	`).Scan(&results)

	c.JSON(http.StatusOK, gin.H{"data": results})
}

/*
========================================
WEEKLY REPORT
========================================
*/

func GetWeeklyReport(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var results []struct {
		Label string  `json:"label"`
		Total float64 `json:"total"`
	}

	db.Raw(`
		SELECT YEARWEEK(date) as label, SUM(total) as total
		FROM transactions
		GROUP BY YEARWEEK(date)
		ORDER BY YEARWEEK(date)
	`).Scan(&results)

	c.JSON(http.StatusOK, gin.H{"data": results})
}

/*
========================================
MONTHLY REPORT
========================================
*/

func GetMonthlyReport(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var results []struct {
		Label string  `json:"label"`
		Total float64 `json:"total"`
	}

	db.Raw(`
		SELECT DATE_FORMAT(date, '%Y-%m') as label, SUM(total) as total
		FROM transactions
		GROUP BY DATE_FORMAT(date, '%Y-%m')
		ORDER BY DATE_FORMAT(date, '%Y-%m')
	`).Scan(&results)

	c.JSON(http.StatusOK, gin.H{"data": results})
}

func DeleteTransaction(c *gin.Context) {
	id := c.Param("id")

	var transaction models.Transaction

	if err := config.DB.Preload("Details").
		First(&transaction, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"message": "Transaksi tidak ditemukan",
		})
		return
	}

	tx := config.DB.Begin()

	if err := tx.Where("transaction_id = ?", transaction.ID).
		Delete(&models.TransactionDetail{}).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal hapus detail transaksi",
		})
		return
	}

	if err := tx.Delete(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal hapus transaksi",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaksi berhasil dihapus",
	})
}

func DeleteAllTransactions(c *gin.Context) {

	tx := config.DB.Begin()

	if err := tx.Exec("DELETE FROM transaction_details").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal hapus detail",
		})
		return
	}

	if err := tx.Exec("DELETE FROM transactions").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal hapus transaksi",
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Semua transaksi berhasil dihapus",
	})
}

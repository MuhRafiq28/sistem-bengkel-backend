package main

import (
	"net/http"
	"backend/controllers"
	"backend/config"
	"backend/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Koneksi database
	config.ConnectDB()

	// Auto migrate
	config.DB.AutoMigrate(
		&models.Product{},
		&models.Transaction{},
		&models.TransactionDetail{},
	)

	// Middleware: set DB ke context
	r.Use(func(c *gin.Context) {
		c.Set("db", config.DB)
		c.Next()
	})

	r.Use(cors.Default())

	api := r.Group("/api")
	{

		// report
		api.GET("/transactions/report/daily", controllers.GetDailyReport)
		api.GET("/transactions/report/weekly", controllers.GetWeeklyReport)
		api.GET("/transactions/report/monthly", controllers.GetMonthlyReport)

		// Product
		api.GET("/products", controllers.GetProducts)
		api.POST("/products", controllers.CreateProduct)
		api.PUT("/products/:id", controllers.UpdateProduct)
		api.DELETE("/products/:id", controllers.DeleteProduct)

		// Transaction
		api.POST("/transactions", controllers.CreateTransaction)
		api.GET("/transactions", controllers.GetTransactionsHistory)

		//history
		api.DELETE("/transactions/:id", controllers.DeleteTransaction)
		api.DELETE("/transactions", controllers.DeleteAllTransactions)

	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Backend Running 🚀"})
	})

	r.Run(":8080")
}

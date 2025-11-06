package routes

import (
	"database/sql"

	"store_product/handlers"
	"store_product/repositories"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes with MySQL
func SetupRoutes(router *gin.Engine, db *sql.DB) {
	// Initialize repositories
	productRepo := repositories.NewProductRepository()
	cartRepo := repositories.NewMySQLCartRepository(db)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	productHandler := handlers.NewProductHandler(productRepo)
	cartHandler := handlers.NewCartHandler(cartRepo)

	setupCommonRoutes(router, healthHandler, productHandler, cartHandler)
}

// SetupRoutesWithDynamoDB configures all application routes with DynamoDB
func SetupRoutesWithDynamoDB(router *gin.Engine, client *dynamodb.Client, tableName string) {
	// Initialize repositories
	productRepo := repositories.NewProductRepository()
	cartRepo := repositories.NewDynamoDBCartRepository(client, tableName)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	productHandler := handlers.NewProductHandler(productRepo)
	cartHandler := handlers.NewCartHandler(cartRepo)

	setupCommonRoutes(router, healthHandler, productHandler, cartHandler)
}

// setupCommonRoutes sets up routes common to all database types
func setupCommonRoutes(router *gin.Engine, healthHandler *handlers.HealthHandler, productHandler *handlers.ProductHandler, cartHandler *handlers.CartHandler) {
	// Health check
	router.GET("/health", healthHandler.Check)

	// Product routes
	router.GET("/products/:productId", productHandler.GetByID)
	router.POST("/products", productHandler.Create)

	// Shopping cart routes
	router.POST("/shopping-carts", cartHandler.Create)
	router.GET("/shopping-carts/:id", cartHandler.GetByID)
	router.POST("/shopping-carts/:id/items", cartHandler.AddItem)
}

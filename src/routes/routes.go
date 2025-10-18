package routes

import (
	"database/sql"

	"store_product/handlers"
	"store_product/repositories"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, db *sql.DB) {
	// Initialize repositories
	productRepo := repositories.NewProductRepository()
	cartRepo := repositories.NewCartRepository(db)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	productHandler := handlers.NewProductHandler(productRepo)
	cartHandler := handlers.NewCartHandler(cartRepo)

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

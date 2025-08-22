package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Product represents the product schema from the OpenAPI spec.
type Product struct {
	ProductID    int    `json:"product_id"`
	SKU          string `json:"sku"`
	Manufacturer string `json:"manufacturer"`
	CategoryID   int    `json:"category_id"`
	Weight       int    `json:"weight"`
	SomeOtherID  int    `json:"some_other_id"`
}

// In-memory product store (thread-safe)
var (
	productStore = make(map[int]Product)
	productMutex sync.RWMutex
)

func main() {
	router := gin.Default()

	// Product endpoints
	router.GET("/products/:productId", getProductByID)
	router.POST("/products", postProduct)

	// Health check endpoint
	router.GET("/health", healthHandler)

	router.Run(":8080")
}

// healthHandler handles GET /health
func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"service":   "product-service",
	})
}

// getProductByID handles GET /products/:productId
func getProductByID(c *gin.Context) {
	idStr := c.Param("productId")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_INPUT",
			"message": "The provided input data is invalid",
			"details": "Product ID must be a positive integer",
		})
		return
	}
	productMutex.RLock()
	product, ok := productStore[id]
	productMutex.RUnlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "INVALID_INPUT",
			"message": "The provided input data is invalid",
			"details": "Product ID must be a positive integer",
		})
		return
	}
	c.JSON(http.StatusOK, product)
}

// postProduct handles POST /products
func postProduct(c *gin.Context) {
	var prod Product
	if err := c.ShouldBindJSON(&prod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_INPUT",
			"message": "The provided input data is invalid",
			"details": "Product ID must be a positive integer",
		})
		return
	}
	if prod.ProductID < 1 || prod.SKU == "" || prod.Manufacturer == "" ||
		prod.CategoryID < 1 || prod.Weight < 0 || prod.SomeOtherID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "INVALID_INPUT",
			"message": "The provided input data is invalid",
			"details": "Product ID must be a positive integer",
		})
		return
	}
	productMutex.Lock()
	productStore[prod.ProductID] = prod
	productMutex.Unlock()
	c.JSON(http.StatusCreated, prod)
}

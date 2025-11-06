package handlers

import (
	"log"
	"net/http"
	"strconv"

	"store_product/models"
	"store_product/repositories"

	"github.com/gin-gonic/gin"
)

// CartHandler handles shopping cart requests
type CartHandler struct {
	repo repositories.CartRepositoryInterface
}

// NewCartHandler creates a new cart handler
func NewCartHandler(repo repositories.CartRepositoryInterface) *CartHandler {
	return &CartHandler{repo: repo}
}

// Create handles POST /shopping-carts
func (h *CartHandler) Create(c *gin.Context) {
	var req models.CreateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	cartID, err := h.repo.Create(req.CustomerID)
	if err != nil {
		log.Printf("Error creating cart: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "DATABASE_ERROR",
			Message: "Failed to create cart",
		})
		return
	}

	// Handle both int (MySQL) and string (DynamoDB) cart IDs
	var response interface{}
	switch v := cartID.(type) {
	case int:
		response = models.CreateCartResponse{ShoppingCartID: v}
	case string:
		response = map[string]string{"shopping_cart_id": v}
	default:
		log.Printf("Unexpected cart ID type: %T", cartID)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "INTERNAL_ERROR",
			Message: "Unexpected cart ID type",
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetByID handles GET /shopping-carts/:id
func (h *CartHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	// Try to parse as int first (MySQL), if it fails, treat as string (DynamoDB UUID)
	var cartID interface{}
	if id, err := strconv.Atoi(idStr); err == nil && id > 0 {
		cartID = id
	} else {
		// Assume it's a UUID string for DynamoDB
		cartID = idStr
	}

	cart, err := h.repo.GetByID(cartID)
	if err != nil {
		log.Printf("Error fetching cart: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "DATABASE_ERROR",
			Message: "Failed to fetch cart",
		})
		return
	}

	if cart == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "NOT_FOUND",
			Message: "Shopping cart not found",
		})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// AddItem handles POST /shopping-carts/:id/items
func (h *CartHandler) AddItem(c *gin.Context) {
	idStr := c.Param("id")

	// Try to parse as int first (MySQL), if it fails, treat as string (DynamoDB UUID)
	var cartID interface{}
	if id, err := strconv.Atoi(idStr); err == nil && id > 0 {
		cartID = id
	} else {
		// Assume it's a UUID string for DynamoDB
		cartID = idStr
	}

	var req models.AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Check if cart exists
	exists, err := h.repo.Exists(cartID)
	if err != nil {
		log.Printf("Error checking cart existence: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "DATABASE_ERROR",
			Message: "Failed to verify cart",
		})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "NOT_FOUND",
			Message: "Shopping cart not found",
		})
		return
	}

	// Add item to cart
	if err := h.repo.AddItem(cartID, req.ProductID, req.Quantity); err != nil {
		log.Printf("Error adding item to cart: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "DATABASE_ERROR",
			Message: "Failed to add item to cart",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

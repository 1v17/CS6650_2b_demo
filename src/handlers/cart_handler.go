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
	repo *repositories.CartRepository
}

// NewCartHandler creates a new cart handler
func NewCartHandler(repo *repositories.CartRepository) *CartHandler {
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

	c.JSON(http.StatusCreated, models.CreateCartResponse{
		ShoppingCartID: cartID,
	})
}

// GetByID handles GET /shopping-carts/:id
func (h *CartHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	cartID, err := strconv.Atoi(idStr)
	if err != nil || cartID <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Invalid cart ID",
			Details: "Cart ID must be a positive integer",
		})
		return
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
	cartID, err := strconv.Atoi(idStr)
	if err != nil || cartID <= 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Invalid cart ID",
			Details: "Cart ID must be a positive integer",
		})
		return
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

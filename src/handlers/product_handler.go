package handlers

import (
	"net/http"
	"strconv"

	"store_product/models"
	"store_product/repositories"

	"github.com/gin-gonic/gin"
)

// ProductHandler handles product-related requests
type ProductHandler struct {
	repo *repositories.ProductRepository
}

// NewProductHandler creates a new product handler
func NewProductHandler(repo *repositories.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// GetByID handles GET /products/:productId
func (h *ProductHandler) GetByID(c *gin.Context) {
	idStr := c.Param("productId")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "The provided input data is invalid",
			Details: "Product ID must be a positive integer",
		})
		return
	}

	product, ok := h.repo.GetByID(id)
	if !ok {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "NOT_FOUND",
			Message: "Product not found",
			Details: "The requested product does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Create handles POST /products
func (h *ProductHandler) Create(c *gin.Context) {
	var prod models.Product
	if err := c.ShouldBindJSON(&prod); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "The provided input data is invalid",
			Details: err.Error(),
		})
		return
	}

	h.repo.Save(prod)
	c.JSON(http.StatusCreated, prod)
}

package models

import "time"

// ShoppingCart represents a customer's shopping cart
type ShoppingCart struct {
	CartID     interface{} `json:"cart_id" dynamodbav:"cart_id"` // Can be int (MySQL) or string (DynamoDB UUID)
	CustomerID int         `json:"customer_id" dynamodbav:"customer_id"`
	CreatedAt  time.Time   `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at" dynamodbav:"updated_at"`
	Items      []CartItem  `json:"items" dynamodbav:"cart_items"`
	TTL        *int64      `json:"ttl,omitempty" dynamodbav:"ttl,omitempty"` // TTL for DynamoDB (Unix timestamp)
}

// CartItem represents an item in a shopping cart
type CartItem struct {
	ItemID    int       `json:"item_id" dynamodbav:"item_id"`
	ProductID int       `json:"product_id" dynamodbav:"product_id"`
	Quantity  int       `json:"quantity" dynamodbav:"quantity"`
	AddedAt   time.Time `json:"added_at" dynamodbav:"added_at"`
	UpdatedAt time.Time `json:"updated_at" dynamodbav:"updated_at"`
}

// CreateCartRequest represents the request body for creating a cart
type CreateCartRequest struct {
	CustomerID int `json:"customer_id" binding:"required,min=1"`
}

// CreateCartResponse represents the response after creating a cart
type CreateCartResponse struct {
	ShoppingCartID int `json:"shopping_cart_id"`
}

// AddItemRequest represents the request body for adding items to cart
type AddItemRequest struct {
	ProductID int `json:"product_id" binding:"required,min=1"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

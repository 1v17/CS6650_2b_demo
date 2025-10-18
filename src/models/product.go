package models

// Product represents the product schema from the OpenAPI spec
type Product struct {
	ProductID    int    `json:"product_id" binding:"required,min=1"`
	SKU          string `json:"sku" binding:"required"`
	Manufacturer string `json:"manufacturer" binding:"required"`
	CategoryID   int    `json:"category_id" binding:"required,min=1"`
	Weight       int    `json:"weight" binding:"required,min=0"`
	SomeOtherID  int    `json:"some_other_id" binding:"required,min=1"`
}

package repositories

import "store_product/models"

// CartRepositoryInterface defines the contract for cart data operations
type CartRepositoryInterface interface {
	Create(customerID int) (interface{}, error)
	GetByID(cartID interface{}) (*models.ShoppingCart, error)
	Exists(cartID interface{}) (bool, error)
	AddItem(cartID interface{}, productID, quantity int) error
	GetByCustomerID(customerID int) ([]models.ShoppingCart, error)
}

package repositories

import (
	"store_product/models"
	"sync"
)

// ProductRepository handles product data operations
type ProductRepository struct {
	store sync.Map
}

// NewProductRepository creates a new product repository
func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(id int) (*models.Product, bool) {
	value, ok := r.store.Load(id)
	if !ok {
		return nil, false
	}
	product := value.(models.Product)
	return &product, true
}

// Save stores a product
func (r *ProductRepository) Save(product models.Product) {
	r.store.Store(product.ProductID, product)
}

// Delete removes a product
func (r *ProductRepository) Delete(id int) {
	r.store.Delete(id)
}

// Exists checks if a product exists
func (r *ProductRepository) Exists(id int) bool {
	_, ok := r.store.Load(id)
	return ok
}

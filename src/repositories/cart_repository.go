package repositories

import (
	"database/sql"
	"fmt"
	"store_product/models"
)

// CartRepository handles shopping cart data operations
type CartRepository struct {
	db *sql.DB
}

// NewCartRepository creates a new cart repository
func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

// Create creates a new shopping cart
func (r *CartRepository) Create(customerID int) (int, error) {
	result, err := r.db.Exec(
		"INSERT INTO shopping_carts (customer_id) VALUES (?)",
		customerID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create cart: %w", err)
	}

	cartID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get cart ID: %w", err)
	}

	return int(cartID), nil
}

// GetByID retrieves a shopping cart by ID with all items
func (r *CartRepository) GetByID(cartID int) (*models.ShoppingCart, error) {
	// Use LEFT JOIN to get cart and all items in a single query
	rows, err := r.db.Query(`
		SELECT 
			c.cart_id, c.customer_id, c.created_at, c.updated_at,
			ci.item_id, ci.product_id, ci.quantity, ci.added_at, ci.updated_at
		FROM shopping_carts c
		LEFT JOIN cart_items ci ON c.cart_id = ci.cart_id
		WHERE c.cart_id = ?
		ORDER BY ci.added_at
	`, cartID)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart: %w", err)
	}
	defer rows.Close()

	var cart *models.ShoppingCart
	cart = nil

	for rows.Next() {
		// Use sql.Null types for the cart_items fields since they may be NULL (LEFT JOIN)
		var itemID, productID, quantity sql.NullInt64
		var addedAt, updatedAt sql.NullTime

		if cart == nil {
			// First row - initialize the cart
			cart = &models.ShoppingCart{Items: []models.CartItem{}}
			err := rows.Scan(
				&cart.CartID, &cart.CustomerID, &cart.CreatedAt, &cart.UpdatedAt,
				&itemID, &productID, &quantity, &addedAt, &updatedAt,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to scan cart: %w", err)
			}
		} else {
			// Subsequent rows - only scan item fields (cart fields are the same)
			var tempCartID, tempCustomerID int
			var tempCreatedAt, tempUpdatedAt interface{}
			err := rows.Scan(
				&tempCartID, &tempCustomerID, &tempCreatedAt, &tempUpdatedAt,
				&itemID, &productID, &quantity, &addedAt, &updatedAt,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to scan cart item: %w", err)
			}
		}

		// Add item to cart if it exists (not NULL from LEFT JOIN)
		if itemID.Valid {
			item := models.CartItem{
				ItemID:    int(itemID.Int64),
				ProductID: int(productID.Int64),
				Quantity:  int(quantity.Int64),
				AddedAt:   addedAt.Time,
				UpdatedAt: updatedAt.Time,
			}
			cart.Items = append(cart.Items, item)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cart rows: %w", err)
	}

	return cart, nil
}

// Exists checks if a cart exists
func (r *CartRepository) Exists(cartID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM shopping_carts WHERE cart_id = ?)",
		cartID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check cart existence: %w", err)
	}
	return exists, nil
}

// AddItem adds or updates an item in the cart
func (r *CartRepository) AddItem(cartID, productID, quantity int) error {
	// Use INSERT ... ON DUPLICATE KEY UPDATE for upsert behavior
	_, err := r.db.Exec(`
		INSERT INTO cart_items (cart_id, product_id, quantity)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			quantity = quantity + VALUES(quantity),
			updated_at = CURRENT_TIMESTAMP
	`, cartID, productID, quantity)

	if err != nil {
		return fmt.Errorf("failed to add item to cart: %w", err)
	}

	// Update cart's updated_at timestamp
	_, err = r.db.Exec(
		"UPDATE shopping_carts SET updated_at = CURRENT_TIMESTAMP WHERE cart_id = ?",
		cartID,
	)
	if err != nil {
		return fmt.Errorf("failed to update cart timestamp: %w", err)
	}

	return nil
}

// GetByCustomerID retrieves all carts for a customer
func (r *CartRepository) GetByCustomerID(customerID int) ([]models.ShoppingCart, error) {
	rows, err := r.db.Query(
		"SELECT cart_id, customer_id, created_at, updated_at FROM shopping_carts WHERE customer_id = ? ORDER BY created_at DESC",
		customerID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch customer carts: %w", err)
	}
	defer rows.Close()

	carts := []models.ShoppingCart{}
	for rows.Next() {
		var cart models.ShoppingCart
		if err := rows.Scan(&cart.CartID, &cart.CustomerID, &cart.CreatedAt, &cart.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan cart: %w", err)
		}
		carts = append(carts, cart)
	}

	return carts, nil
}

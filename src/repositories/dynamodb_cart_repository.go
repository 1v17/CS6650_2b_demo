package repositories

import (
	"context"
	"fmt"
	"time"

	"store_product/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

// DynamoDBCartRepository handles shopping cart data operations for DynamoDB
type DynamoDBCartRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoDBCartRepository creates a new DynamoDB cart repository
func NewDynamoDBCartRepository(client *dynamodb.Client, tableName string) *DynamoDBCartRepository {
	return &DynamoDBCartRepository{
		client:    client,
		tableName: tableName,
	}
}

// Ensure DynamoDBCartRepository implements CartRepositoryInterface
var _ CartRepositoryInterface = (*DynamoDBCartRepository)(nil)

// Create creates a new shopping cart
func (r *DynamoDBCartRepository) Create(customerID int) (interface{}, error) {
	cartID := uuid.New().String()
	now := time.Now()
	ttl := now.Add(24 * time.Hour).Unix() // 24 hours from now

	cart := models.ShoppingCart{
		CartID:     cartID,
		CustomerID: customerID,
		CreatedAt:  now,
		UpdatedAt:  now,
		Items:      []models.CartItem{},
		TTL:        &ttl,
	}

	// Marshal the cart to DynamoDB attribute values
	av, err := attributevalue.MarshalMap(cart)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cart: %w", err)
	}

	// Put item into DynamoDB
	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cart in DynamoDB: %w", err)
	}

	return cartID, nil
}

// GetByID retrieves a shopping cart by ID with all items
func (r *DynamoDBCartRepository) GetByID(cartID interface{}) (*models.ShoppingCart, error) {
	id, ok := cartID.(string)
	if !ok {
		return nil, fmt.Errorf("invalid cart ID type for DynamoDB")
	}

	// Get item from DynamoDB with strong consistency
	result, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"cart_id": &types.AttributeValueMemberS{Value: id},
		},
		ConsistentRead: aws.Bool(true), // Strong consistency
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cart from DynamoDB: %w", err)
	}

	// Check if item exists
	if result.Item == nil {
		return nil, nil // Cart not found
	}

	// Unmarshal the result into ShoppingCart
	var cart models.ShoppingCart
	err = attributevalue.UnmarshalMap(result.Item, &cart)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart: %w", err)
	}

	return &cart, nil
}

// Exists checks if a cart exists
func (r *DynamoDBCartRepository) Exists(cartID interface{}) (bool, error) {
	cart, err := r.GetByID(cartID)
	if err != nil {
		return false, err
	}
	return cart != nil, nil
}

// AddItem adds or updates an item in the cart
func (r *DynamoDBCartRepository) AddItem(cartID interface{}, productID, quantity int) error {
	id, ok := cartID.(string)
	if !ok {
		return fmt.Errorf("invalid cart ID type for DynamoDB")
	}

	// First, get the current cart
	cart, err := r.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}
	if cart == nil {
		return fmt.Errorf("cart not found")
	}

	// Check if product already exists in items
	found := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			// Increment quantity
			cart.Items[i].Quantity += quantity
			cart.Items[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	// If not found, append new item
	if !found {
		newItem := models.CartItem{
			ItemID:    len(cart.Items) + 1, // Simple incrementing ID
			ProductID: productID,
			Quantity:  quantity,
			AddedAt:   time.Now(),
			UpdatedAt: time.Now(),
		}
		cart.Items = append(cart.Items, newItem)
	}

	// Update the cart's updated_at timestamp
	cart.UpdatedAt = time.Now()

	// Marshal the updated items list
	itemsAV, err := attributevalue.Marshal(cart.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal items: %w", err)
	}

	updatedAtAV, err := attributevalue.Marshal(cart.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to marshal updated_at: %w", err)
	}

	// Update the cart in DynamoDB
	_, err = r.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"cart_id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression: aws.String("SET cart_items = :cart_items, updated_at = :updated_at"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cart_items": itemsAV,
			":updated_at": updatedAtAV,
		},
		ConditionExpression: aws.String("attribute_exists(cart_id)"), // Ensure cart exists
	})
	if err != nil {
		return fmt.Errorf("failed to update cart in DynamoDB: %w", err)
	}

	return nil
}

// GetByCustomerID retrieves all carts for a customer using GSI
func (r *DynamoDBCartRepository) GetByCustomerID(customerID int) ([]models.ShoppingCart, error) {
	// Query using the customer-index GSI
	result, err := r.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("customer-index"),
		KeyConditionExpression: aws.String("customer_id = :customer_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":customer_id": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", customerID)},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query carts by customer ID: %w", err)
	}

	// Unmarshal the results
	var carts []models.ShoppingCart
	err = attributevalue.UnmarshalListOfMaps(result.Items, &carts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal carts: %w", err)
	}

	return carts, nil
}

package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	_ "github.com/go-sql-driver/mysql"
)

// DBConfig holds database configuration
type DBConfig struct {
	Type           string // "mysql" or "dynamodb"
	MySQLDB        *sql.DB
	DynamoDBClient *dynamodb.Client
	DynamoDBTable  string
	AWSRegion      string
}

// GetDatabaseType returns the configured database type
func GetDatabaseType() string {
	dbType := os.Getenv("DATABASE_TYPE")
	if dbType == "" {
		dbType = "mysql" // default to MySQL for backward compatibility
	}
	return dbType
}

func InitDB() (*sql.DB, error) {
	required := []string{"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME"}
	values := make(map[string]string, len(required))
	var missing []string

	for _, k := range required {
		v, ok := os.LookupEnv(k)
		if !ok || v == "" {
			missing = append(missing, k)
			continue
		}
		values[k] = v
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missing)
	}

	// Build DSN from validated environment variables
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		values["DB_USER"],
		values["DB_PASSWORD"],
		values["DB_HOST"],
		values["DB_PORT"],
		values["DB_NAME"],
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Connection pool configuration
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")

	// Initialize schema
	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Println("Database schema initialized")
	return db, nil
}

func initSchema(db *sql.DB) error {
	// Create shopping_carts table
	createCartsTable := `
	CREATE TABLE IF NOT EXISTS shopping_carts (
		cart_id INT AUTO_INCREMENT PRIMARY KEY,
		customer_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_customer_id (customer_id)
	) ENGINE=InnoDB`

	if _, err := db.Exec(createCartsTable); err != nil {
		return fmt.Errorf("failed to create shopping_carts table: %w", err)
	}

	// Create cart_items table
	createItemsTable := `
	CREATE TABLE IF NOT EXISTS cart_items (
		item_id INT AUTO_INCREMENT PRIMARY KEY,
		cart_id INT NOT NULL,
		product_id INT NOT NULL,
		quantity INT NOT NULL DEFAULT 1,
		added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (cart_id) REFERENCES shopping_carts(cart_id) ON DELETE CASCADE,
		UNIQUE KEY unique_cart_product (cart_id, product_id),
		INDEX idx_cart_id (cart_id),
		CHECK (quantity >= 0)
	) ENGINE=InnoDB`

	if _, err := db.Exec(createItemsTable); err != nil {
		return fmt.Errorf("failed to create cart_items table: %w", err)
	}

	log.Println("Database schema initialized successfully")
	return nil
}

// InitDynamoDB initializes DynamoDB client
func InitDynamoDB() (*dynamodb.Client, string, error) {
	// Check required environment variables
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		return nil, "", fmt.Errorf("DYNAMODB_TABLE_NAME environment variable is required")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-west-2" // default region
		log.Printf("AWS_REGION not set, using default: %s", region)
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create DynamoDB client
	client := dynamodb.NewFromConfig(cfg)

	log.Printf("Successfully initialized DynamoDB client for table: %s in region: %s", tableName, region)
	return client, tableName, nil
}

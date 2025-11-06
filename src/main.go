package main

import (
	"log"
	"os"

	"store_product/config"
	"store_product/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables from .env file
	// if err := godotenv.Load(); err != nil {
	// 	log.Printf("Warning: Could not load .env file: %v", err)
	// }

	// Get database type from environment
	dbType := config.GetDatabaseType()
	log.Printf("Database type: %s", dbType)

	// Initialize Gin router
	router := gin.Default()

	// Setup routes based on database type
	if dbType == "dynamodb" {
		// Initialize DynamoDB
		dynamoClient, tableName, err := config.InitDynamoDB()
		if err != nil {
			log.Fatal("Failed to initialize DynamoDB:", err)
		}
		log.Printf("DynamoDB initialized successfully with table: %s", tableName)

		routes.SetupRoutesWithDynamoDB(router, dynamoClient, tableName)
	} else {
		// Initialize MySQL (default)
		db, err := config.InitDB()
		if err != nil {
			log.Fatal("Failed to connect to MySQL database:", err)
		}
		defer db.Close()
		log.Println("MySQL database connection established and successfully initialized.")

		routes.SetupRoutes(router, db)
	}

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

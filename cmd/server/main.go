package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	delivery "polyglot/suppliers/internal/delivery/http"
	"polyglot/suppliers/internal/repository"
	"polyglot/suppliers/internal/usecase"
)

func main() {
	log.Println("Starting Suppliers API Service...")

	// 1. Fetch environment configuration (Fallbacks for local run outside Docker)
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "polyglot_inventory"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 2. Initialize Database Connection Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Connecting to MongoDB at: %s", mongoURI)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB client options: %v", err)
	}

	// Verify the network connection to the database cluster
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Could not ping MongoDB cluster: %v", err)
	}
	log.Println("Successfully connected to MongoDB cluster!")

	db := client.Database(dbName)

	// 3. Clean Architecture Dependency Injection (Bottom to Top)
	repo := repository.NewMongoSupplierRepository(db, "suppliers")
	uCase := usecase.NewSupplierUseCase(repo)

	// 4. Initialize HTTP Delivery Router Engine
	router := gin.Default()

	// Base sanity health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "suppliers-api"})
	})

	// Register our functional domain routes
	delivery.NewSupplierHandler(router, uCase)

	// 5. Start Server Listener Engine
	log.Printf("Suppliers API server running on port :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start web server router: %v", err)
	}
}

package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"suppliers-api/config"
	delivery "suppliers-api/internal/delivery/http"
	"suppliers-api/internal/repository"
	"suppliers-api/internal/usecase"
)

func main() {
	log.Println("Starting Suppliers API Service...")

	// 1. Fetch environment configuration (Fallbacks for local run outside Docker)
	// 1. Initialize Secure Configuration Matrix via HashiCorp Vault
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Printf("Vault initialization warning: %v. Proceeding with structural fallbacks.", err)
		appConfig = &config.AppConfig{
			MongoURI: "mongodb://localhost:27017",
			DBName:   "polyglot_inventory",
			Port:     "8080",
		}
	}

	// 2. Initialize Database Connection Context using safe variables
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Connecting to MongoDB via configuration registry...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(appConfig.MongoURI))

	// Verify the network connection to the database cluster
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Could not ping MongoDB cluster: %v", err)
	}
	log.Println("Successfully connected to MongoDB cluster!")

	db := client.Database(appConfig.DBName)

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
	log.Printf("Suppliers API server running on port :%s", appConfig.Port) // Fix here
	if err := router.Run(":" + appConfig.Port); err != nil {               // Fix here
		log.Fatalf("Failed to start web server router: %v", err)
	}
}

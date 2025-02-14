package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sean-miningah/sil-backend-assessment/internal/adapters/handlers/graphql"
	"github.com/sean-miningah/sil-backend-assessment/internal/adapters/handlers/rest"
	"github.com/sean-miningah/sil-backend-assessment/internal/adapters/notification"
	repo "github.com/sean-miningah/sil-backend-assessment/internal/adapters/repositories/postgres"
	"github.com/sean-miningah/sil-backend-assessment/internal/services"
	"github.com/sean-miningah/sil-backend-assessment/pkg/auth"
	"github.com/sean-miningah/sil-backend-assessment/pkg/auth/middleware"
	"github.com/sean-miningah/sil-backend-assessment/pkg/config"
	"github.com/sean-miningah/sil-backend-assessment/pkg/database"
	"github.com/sean-miningah/sil-backend-assessment/pkg/telemetry"
)

// Application should
// Input and upload products with their various categories
// Return average product price for a category
// Make Order

// Auth using openid connect
// When order is made send customer sms alerting them using Africa's Talking API
// Send Administrator an email about order placed
// Deploy to k8s
// Write Docs

func main() {
	cfg := config.Load(".env")

	//Initialize tracer
	tp, err := telemetry.InitTracer(cfg.ServiceName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer: %v", err)
		}
	}()

	// Construct the PostgreSQL connection string
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	// Initialize database connection
	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repository
	productRepo := repo.NewProductRepository(db)
	orderRepo := repo.NewOrderRepository(db)
	customerRepo := repo.NewCustomerRepoisotory(db)
	notificationRepo := notification.NewNotificationRepo(cfg.ATAPIKey, cfg.NotificationUsername, cfg.GmailAppAPIKey, cfg.ATAPIUrl)

	// Initialize service
	productService := services.NewProductService(productRepo, orderRepo)
	orderService := services.NewOrderService(orderRepo, productRepo, notificationRepo)
	customerService := services.NewCustomerService(customerRepo)

	// Initialize handler
	productHandler := rest.NewProductHandler(productService)
	orderHandler := rest.NewOrderHandler(orderService)

	// Initialize GraphQL handler
	graphqlHandler := graphql.NewHandler(productService, orderService)

	authConfig := auth.NewAuthConfig(cfg)
	authHandler := rest.NewAuthHandler(authConfig, customerService)

	// Gin router setup
	router := gin.Default()

	router.GET("/auth/login", authHandler.Login)
	router.GET("/auth/google/callback", authHandler.GoogleCallback)
	router.POST("/auth/logout", authHandler.Logout)

	// Register routes
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(authConfig.JWTSecret))
	{
		// api.GET("/products", productHandler.List)
		// api.GET("/products/:id", productHandler.Get)
		api.POST("/products", productHandler.Create)
		// api.PUT("/products/:id", productHandler.Update)
		// api.DELETE("/products/:id", productHandler.Delete)

		api.GET("/categories/:categoryId/average-price", productHandler.GetAveragePriceByCategory)

		// Order Routes
		// api.GET("/orders", orderHandler.List)
		// api.GET("/orders/:id", orderHandler.Get)
		api.POST("/orders", orderHandler.Create)
		// api.PUT("/orders/:id", orderHandler.Update)
		// api.DELETE("/order/:id", orderHandler.Delete)
	}

	router.POST("/graphql", graphqlHandler.GraphQL())
	if cfg.Environment == "development" {
		router.GET("/playground", graphqlHandler.Playground())
	}

	// Start the server
	if err := router.Run(cfg.Address); err != nil {
		log.Fatal(err)
	}
}

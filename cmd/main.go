package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	migrate_driver "github.com/golang-migrate/migrate/v4/database/cockroachdb" // migrate_driver "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"                       // Import the file source driver

	"stock-api/config"
	"stock-api/infrastructure"
	"stock-api/infrastructure/adapters/handler"
	"stock-api/infrastructure/adapters/middleware"
	"stock-api/infrastructure/adapters/repository"
	"stock-api/infrastructure/core/domain"
	"stock-api/infrastructure/core/service"
)

var (
	mode         = flag.String("mode", "api", "Mode: 'api' or 'data'")
	migrate_dir  = flag.String("migrate", "", "Run database migrations 'up' or 'down'")
	repo         *repository.StockBDRepository
	stockService *service.StockService
	httpHandler  *handler.StockHandler
)

// Sets up the Gin router with middleware
func setupRouter() *gin.Engine {
	r := gin.Default()

	// Use middlewares
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	return r
}

// Defines the API routes
func setupRoutes(router *gin.Engine) {
	srv := service.NewBestInvestmentsService()
	httpHandler = handler.NewStockHandler(stockService, srv)
	api := router.Group("/api/v1")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	api.POST("/stocks", httpHandler.FindStocks)
	api.GET("/recommendations", httpHandler.GetStockRecommendations)
}

// Runs database migrations in the specified direction ('up' or 'down')
func RunMigrations(cfg *config.Config, db *sql.DB, direction string) error {
	// Validate the direction argument
	if direction != "up" && direction != "down" {
		return fmt.Errorf("invalid migration direction: %s", direction)
	}

	// Run migrations
	driver, err := migrate_driver.WithInstance(db, &migrate_driver.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "migrations"), // Path to migrations
		cfg.DB.DBName,                          // Database name
		driver,
	)
	if err != nil {
		return fmt.Errorf("error initializing migrations: %w", err)
	}

	// Run migrations based on the specified direction
	switch direction {
	case "up":
		if err := m.Up(); err != nil {
			return fmt.Errorf("error applying migrations: %w", err)
		}
		log.Println("Migrations applied successfully")
	case "down":
		if err := m.Down(); err != nil {
			return fmt.Errorf("error rolling back migrations: %w", err)
		}
		log.Println("Migrations rolled back successfully")
	}

	return nil
}

// Sets up the batch processor and runs it in a goroutine
func setupBatchProcessor(cfg *config.Config, done chan struct{}) {
	apiClient := service.NewExternalAPIClient(cfg.ExternalAPI.URL)
	classificationService := service.NewClassificationService()

	processor := handler.NewBatchProcessor(
		apiClient,
		repo,
		classificationService,
		cfg.ExternalAPI.BatchSize,
		cfg.ExternalAPI.JWTToken,
		500, // e.g., 500ms
	)

	go func() {
		defer close(done) // Closes the channel when the process finishes
		if err := processor.ProcessStocks(context.Background()); err != nil {
			log.Printf("Error processing stocks: %v", err)
		}
	}()
}

func main() {
	flag.Parse()
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize the database connection
	db, err := infrastructure.NewDatabaseConnection(cfg.DB)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return // Ensure deferred functions are executed
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Error getting database instance:", err)
		return // Ensure deferred functions are executed
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()
	log.Println("Database connection established")

	// Initialize the repository
	repo = repository.NewStockBDRepository(db)
	if repo == nil {
		log.Println("Error initializing repository")
		return
	}
	log.Println("Repository initialized")

	// Initialize the service
	stockService = service.NewStockService(repo, repository.NewGormFieldValidator(&domain.Stock{}))
	if stockService == nil {
		log.Println("Error initializing service")
		return
	}
	log.Println("Service initialized")

	if *migrate_dir != "" {
		// Run database migrations
		log.Printf("Running migrations: %s", *migrate_dir)

		if err := RunMigrations(cfg, sqlDB, *migrate_dir); err != nil {
			log.Printf("Error running migrations: %v", err)
			return
		}
		log.Println("Migrations completed")
		return
	}

	switch *mode {
	case "api":
		// Setting up the Gin router
		router := setupRouter()

		// Setting up the routes
		setupRoutes(router)

		// HTTP Server with graceful shutdown
		srv := &http.Server{
			Addr:              fmt.Sprintf("%s:%d", cfg.Server.URL, cfg.Server.Port),
			Handler:           router,
			ReadHeaderTimeout: 10 * time.Second, // Add a timeout for reading headers
		}

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()
		log.Printf("Server started on port %d", cfg.Server.Port)
	case "data":
		// Setting up the batch processor
		done := make(chan struct{}) // Channel to coordinate shutdown
		setupBatchProcessor(cfg, done)
		log.Println("Batch processor started")

		// Wait for the goroutine to finish
		<-done
		log.Println("Batch processor finished")
	default:
		return
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}

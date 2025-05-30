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
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	migrate_driver "github.com/golang-migrate/migrate/v4/database/cockroachdb" // migrate_driver "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"                       // Import the file source driver
	"go.uber.org/zap"

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

// setupRouter configures the Gin router with all required middleware.
// It sets up CORS, logging, and recovery middleware.
// Returns a configured *gin.Engine instance.
func setupRouter(cfg *config.Config, zapLogger *zap.Logger) *gin.Engine {
	r := gin.Default()

	// Register middlewares
	r.Use(middleware.AsyncCORSMiddleware(cfg.AllowedOrigins))
	r.Use(middleware.AsyncLogger(zapLogger))
	r.Use(gin.Recovery())

	return r
}

// setupRoutes defines all API endpoints and attaches them to the router.
// It initializes the handler with the worker pool and services.
func setupRoutes(router *gin.Engine) {
	srv := service.NewBestInvestmentsService()

	// Worker pool size = (cores * 2) + 1 (for storage units)
	workerPoolSize := (runtime.NumCPU() * 2) + 1

	httpHandler = handler.NewStockHandler(stockService, srv, workerPoolSize)
	api := router.Group("/api/v1")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	api.POST("/stocks", httpHandler.FindStocks)
	api.GET("/recommendations", httpHandler.GetStockRecommendations)
}

// RunMigrations executes database migrations in the specified direction ("up" or "down").
// It initializes the migration driver and runs the migrations from the "migrations" directory.
// Returns an error if migration fails.
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

// setupBatchProcessor initializes and runs the batch processor in a goroutine.
// It processes stocks using the external API client and classification service.
// The done channel is closed when processing is finished.
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

// main is the entry point of the application.
// It loads configuration, initializes the database, repository, and services,
// and starts the API server or batch processor based on the selected mode.
// It also handles graceful shutdown on interrupt signals.
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
		zapLogger, err := zap.NewProduction()
		if err != nil {
			log.Printf("Failed to initialize zap logger: %v", err)
			return
		}
		defer func() {
			if err := zapLogger.Sync(); err != nil && !strings.Contains(err.Error(), "invalid argument") {
				log.Printf("Error syncing zap logger: %v", err)
			}
		}()

		router := setupRouter(cfg, zapLogger)

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

package infrastructure

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"stock-api/config"
)

// NewDatabaseConnection creates a new database connection for PostgreSQL or CockroachDB
func NewDatabaseConnection(cfg config.DBConfig) (*gorm.DB, error) {
	var dsn string

	// Build DSN based on the database type
	switch cfg.DBType {
	case "postgresql":
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode, cfg.TimeZone,
		)
	case "cockroachdb":
		dsn = fmt.Sprintf(
			"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
		)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}

	// Open the database connection
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true, // Improves performance with prepared statements
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %w", err)
	}

	// Configure the connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Printf("Successfully connected to %s database", cfg.DBType)
	return db, nil
}

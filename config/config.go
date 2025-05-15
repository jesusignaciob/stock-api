package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// ExternalAPIConfig holds the configuration for an external API.
// Fields:
// - URL: The base URL of the external API.
// - JWTToken: The JWT token used for authentication with the external API.
// - BatchSize: The size of batches for API requests.
type ExternalAPIConfig struct {
	URL       string
	JWTToken  string
	BatchSize int
}

// ServerConfig holds the configuration for the server.
// Fields:
// - URL: The base URL of the server.
// - Port: The port on which the server listens.
type ServerConfig struct {
	URL  string
	Port int
}

// DBConfig holds the configuration for the database connection.
// Fields:
// - DBType: The type of database (e.g., PostgreSQL, CockroachDB).
// - Host: The hostname or IP address of the database server.
// - Port: The port on which the database server listens.
// - User: The username for database authentication.
// - Password: The password for database authentication.
// - DBName: The name of the database to connect to.
// - SSLMode: The SSL mode for the database connection (e.g., "disable", "require").
// - TimeZone: The timezone for the database connection.
type DBConfig struct {
	DBType   string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

// Config holds the overall application configuration.
// Fields:
// - ExternalAPI: Configuration for the external API.
// - Server: Configuration for the server.
// - DB: Configuration for the database.
type Config struct {
	ExternalAPI ExternalAPIConfig
	Server      ServerConfig
	DB          DBConfig
}

// LoadConfig loads the application configuration from environment variables or a .env file.
// It initializes the configuration for the external API, server, and database.
//
// Returns:
// - A pointer to a Config struct containing the loaded configuration.
// - An error if any required configuration value cannot be parsed.
func LoadConfig() (*Config, error) {
	// Load the .env file if it exists.
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Parse the batch size for the external API.
	batchSize, err := strconv.Atoi(getEnv("EXTERNAL_API_BATCH_SIZE", "100"))
	if err != nil {
		return nil, err
	}

	// Parse the server port.
	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, err
	}

	// Parse the database port.
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, err
	}

	// Initialize the configuration struct.
	cfg := &Config{
		ExternalAPI: ExternalAPIConfig{
			URL:       getEnv("EXTERNAL_API_URL", "https://api.example.com"),
			JWTToken:  getEnv("EXTERNAL_API_JWT_TOKEN", "your_jwt_token"),
			BatchSize: batchSize,
		},
		Server: ServerConfig{
			URL:  getEnv("SERVER_URL", "https://app.example.com"),
			Port: port,
		},
		DB: DBConfig{
			DBType:   getEnv("DB_TYPE", "cockroachdb"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "api_user"),
			Password: getEnv("DB_PASSWORD", "P@ssw0rd"),
			DBName:   getEnv("DB_NAME", "api_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			TimeZone: "UTC",
		},
	}

	return cfg, nil
}

// getEnv retrieves the value of an environment variable.
// If the variable is not set, it returns the provided default value.
//
// Parameters:
// - key: The name of the environment variable.
// - defaultValue: The default value to return if the variable is not set.
//
// Returns:
// - The value of the environment variable, or the default value if the variable is not set.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

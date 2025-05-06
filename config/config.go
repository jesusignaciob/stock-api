package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ExternalAPIConfig struct {
	URL       string
	JWTToken  string
	BatchSize int
}

type ServerConfig struct {
	URL  string
	Port int
}

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

type Config struct {
	ExternalAPI ExternalAPIConfig
	Server      ServerConfig
	DB          DBConfig
}

func LoadConfig() (*Config, error) {
	// Cargar .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Configuración de la API externa
	batchSize, err := strconv.Atoi(getEnv("EXTERNAL_API_BATCH_SIZE", "100"))
	if err != nil {
		return nil, err
	}

	// Configuración del servidor
	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, err
	}

	// Configuración de la base de datos
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, err
	}

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
			TimeZone: getEnv("DB_TIMEZONE", "UTC"),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

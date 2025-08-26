package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI    string
	DatabaseName string
	Environment string
	Port        string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	env := getEnv("GO_ENV", "production")
	
	// Determine database name based on environment
	var dbName string
	switch env {
	case "development", "test":
		dbName = "GoShopping-test"
	case "production":
		dbName = "GoShopping"
	default:
		log.Printf("Warning: Unknown environment '%s', defaulting to production", env)
		dbName = "GoShopping"
		env = "production"
	}

	mongoURI := os.Getenv("GO_SHOPPING_MONGO_ATLAS_URI")
	if mongoURI == "" {
		log.Fatal("Error: GO_SHOPPING_MONGO_ATLAS_URI environment variable is not set")
	}

	return &Config{
		MongoURI:    mongoURI,
		DatabaseName: dbName,
		Environment: env,
		Port:        getEnv("PORT", "8080"),
	}
}

// getEnv gets environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsDevelopment returns true if running in development/test environment
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development" || c.Environment == "test"
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
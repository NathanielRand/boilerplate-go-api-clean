// config/config.go

package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Load loads environment variables from a .env file
func Load() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file:", err)
		return err
	}
	return nil
}

// Get retrieves the value of an environment variable by name
func Get(name string) string {
	return os.Getenv(name)
}

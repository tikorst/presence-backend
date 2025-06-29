package config

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Global variable to hold the database connection
var DB *gorm.DB

// ConnectDB initializes the database connection
func ConnectDB() {
	// Load the DSN from environment variable
	// Ensure that the DB_DSN environment variable is set
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		panic("DB_DSN environment variable is not set")
	}

	// Open a new database connection using GORM
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		TranslateError: true,
	})

	// Check for errors in the connection
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// Check if the database connection is established
	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Print success message
	fmt.Println("Database connected successfully!")
}

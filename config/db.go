package config

// import (
// 	"crypto/tls"
// 	"crypto/x509"
// 	"database/sql"
// 	"fmt"
// 	"os"

// 	"github.com/go-sql-driver/mysql"
// 	_ "github.com/go-sql-driver/mysql"
// 	"gorm.io/gorm"
// )

// var DB *gorm.DB

// func ConnectDB() (*sql.DB, error) {
// 	// Baca certificate SSL
// 	certBytes, err := os.ReadFile("ssl/DigiCertGlobalRootCA.crt.pem")
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Register SSL certificate
// 	rootCAs := x509.NewCertPool()
// 	if ok := rootCAs.AppendCertsFromPEM(certBytes); !ok {
// 		return nil, fmt.Errorf("failed to append CA certificate")
// 	}

// 	mysql.RegisterTLSConfig("custom", &tls.Config{
// 		RootCAs: rootCAs,
// 	})
// 	dsn := os.Getenv("DB_DSN")
// 	DB, err = gorm.Open("mysql", dsn)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }
// database.go

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
	sqlDB.SetMaxOpenConns(130)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Print success message
	fmt.Println("Database connected successfully!")
}

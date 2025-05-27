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

var DB *gorm.DB

// ConnectDB initializes the database connection
func ConnectDB() {
	// Ensure TLS is initialized before connecting
	InitTLS()

	// Define DSN with custom TLS
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		panic("DB_DSN environment variable is not set")
	}
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)           // jumlah koneksi idle maksimum
	sqlDB.SetMaxOpenConns(341)          // jumlah koneksi maksimum yang terbuka
	sqlDB.SetConnMaxLifetime(time.Hour) // batas waktu koneksi aktif
	fmt.Println("Database connected successfully!")
}

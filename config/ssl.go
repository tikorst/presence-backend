package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

// InitTLS registers a custom TLS config for MySQL
func InitTLS() {
	rootCertPool := x509.NewCertPool()

	// Read CA certificate
	caCert, err := os.ReadFile("ssl/DigiCertGlobalRootCA.crt.pem")
	if err != nil {
		fmt.Println("Failed to read CA cert:", err)
		return
	}
	rootCertPool.AppendCertsFromPEM(caCert)

	// Register custom TLS config
	err = mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs: rootCertPool,
	})
	if err != nil {
		fmt.Println("Failed to register TLS config:", err)
	}
}

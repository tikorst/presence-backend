package middleware

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
)

// AppCheckMiddleware memverifikasi token Firebase App Check dari header.
func AppCheckMiddleware(app *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari header
		appCheckToken := c.GetHeader("X-Firebase-AppCheck")
		if appCheckToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: App Check token not found."})
			c.Abort() // Hentikan request
			return
		}

		// Buat client App Check dari instance Firebase App Anda
		client, err := app.AppCheck(context.Background())
		if err != nil {
			log.Printf("Error creating App Check client: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
			c.Abort()
			return
		}

		// Verifikasi token. Ini bisa memverifikasi token asli maupun token debug.
		_, err = client.VerifyToken(appCheckToken)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Invalid App Check token.", "details": err.Error()})
			c.Abort()
			return
		}

		// Sukses! Lanjutkan ke handler utama.
		c.Next()
	}
}

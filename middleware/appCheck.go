package middleware

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
)

// AppCheckMiddleware verifies Firebase App Check tokens in incoming requests.
func AppCheckMiddleware(app *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the App Check token from the request header
		appCheckToken := c.GetHeader("X-Firebase-AppCheck")
		if appCheckToken == "" {

			// if the token is not present, return an error response
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: App Check token not found."})
			c.Abort()
			return
		}

		// Create an App Check client using the Firebase app instance
		client, err := app.AppCheck(context.Background())
		if err != nil {

			// Log the error and return an internal server error response
			log.Printf("Error creating App Check client: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
			c.Abort()
			return
		}

		// Verify the App Check token
		_, err = client.VerifyToken(appCheckToken)
		if err != nil {

			// If verification fails, log the error and return a forbidden response
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Invalid App Check token.", "details": err.Error()})
			c.Abort()
			return
		}

		// If verification is successful, log the success and proceed to the next handler
		c.Next()
	}
}

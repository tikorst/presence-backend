package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Auth is a middleware function that checks for a valid JWT token in the request header.
func Auth() gin.HandlerFunc {

	// return a middleware function that checks the JWT token
	return func(c *gin.Context) {

		// Get the token from the Authorization header
		tokenString := c.GetHeader("Authorization")

		// If the token is empty, return an unauthorized response
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			c.Abort()
			return
		}

		// If the token starts with "Bearer ", remove it
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// Parse the token using the jwt package
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret key used to sign the token
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		// If there was an error parsing the token, return an unauthorized response
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// If the token is valid, check the claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Expired token"})
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			// Set the claims in the context for further use
			c.Set("claims", claims)
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

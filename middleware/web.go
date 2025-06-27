package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Web is a middleware function that checks for a valid JWT token in the request cookie.
func Web() gin.HandlerFunc {

	// Return a middleware function that checks the JWT token
	return func(c *gin.Context) {

		// Get the token from the cookie
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			c.Abort()
			return
		}

		// If the token starts with "Bearer ", remove it
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			fmt.Printf("Bearer token detected: %s\n", tokenString)
			tokenString = tokenString[7:]
		}

		// Parse the token using the jwt package
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SECRET_KEY")), nil
		})


		// If there was an error parsing the token, return an unauthorized response
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Check if the token is valid
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Check if the token is expired
		if exp, ok := claims["exp"].(float64); !ok || float64(time.Now().Unix()) > exp {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Expired token"})
			return
		}

		// Set the claims in the context for further use
		c.Set("claims", claims)
		c.Next()
	}
}

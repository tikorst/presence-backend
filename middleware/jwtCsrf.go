package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tikorst/presence-backend/helpers"
)

// Generate CSRF token

// JWT CSRF Middleware
func JWTCSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// skip CSRF check for GET and OPTIONS requests
		if c.Request.Method == "GET" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}


		// Get the JWT token from the cookie
		jwtCSRFToken,_ := helpers.GetCSRFToken(c)
		if jwtCSRFToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No CSRF token in JWT"})
			return
		}

		// Get the CSRF token from the request header
		headerCSRFToken := c.GetHeader("X-CSRF-Token")
		if headerCSRFToken == "" || headerCSRFToken != jwtCSRFToken {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF token mismatch"})
			return
		}

		// Parse the JWT token
		c.Next()
	}
}

// Extract CSRF token from JWT
func GetCSRFFromJWT(c *gin.Context) {

	// Get the JWT token from the cookie
	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	// Parse the JWT token
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	// Check if the token is valid
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Get the CSRF token from the claims
	csrfToken, exists := claims["csrf_token"].(string)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "No CSRF token found"})
		return
	}

	// Set the CSRF token in the response
	c.JSON(http.StatusOK, gin.H{
		"csrf_token": csrfToken,
		"expires_in": 3600,
	})
}

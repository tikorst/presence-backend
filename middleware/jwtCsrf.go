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
		// Skip untuk GET dan OPTIONS
		if c.Request.Method == "GET" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		jwtCSRFToken,_ := helpers.GetCSRFToken(c)
		if jwtCSRFToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No CSRF token in JWT"})
			return
		}

		headerCSRFToken := c.GetHeader("X-CSRF-Token")
		if headerCSRFToken == "" || headerCSRFToken != jwtCSRFToken {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF token mismatch"})
			return
		}

		c.Next()
	}
}

// Extract CSRF token dari JWT
func GetCSRFFromJWT(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	csrfToken, exists := claims["csrf_token"].(string)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "No CSRF token found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"csrf_token": csrfToken,
		"expires_in": 3600,
	})
}

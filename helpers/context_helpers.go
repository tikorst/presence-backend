package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Ambil semua claims dari context
func GetClaims(c *gin.Context) (jwt.MapClaims, error) {
	claimsRaw, exists := c.Get("claims")
	if !exists {
		return nil, errors.New("jwt claims not found in context")
	}

	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid jwt claims format")
	}

	return claims, nil
}

func getClaim(c *gin.Context, key string) (string, error) {
	claims, err := GetClaims(c)
	if err != nil {
		return "", err
	}

	value, ok := claims[key].(string)
	if !ok {
		return "", errors.New(key + " claim not found or invalid")
	}

	return value, nil
}

// Ambil username dari sub claim
func GetUsername(c *gin.Context) (string, error) {
	return getClaim(c, "sub")
}

// Ambil role user dari claim
func GetRole(c *gin.Context) (string, error) {
	return getClaim(c, "role")
}

// Ambil device ID dari claim
func GetDeviceID(c *gin.Context) (string, error) {
	return getClaim(c, "device_id")
}

// Ambil CSRF token dari context
func GetCSRFToken(c *gin.Context) (string, error) {
	return getClaim(c, "csrf_token")
}

package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Get all JWT claims from the context
func GetClaims(c *gin.Context) (jwt.MapClaims, error) {
	// Check if the claims are present in the context
	claimsRaw, exists := c.Get("claims")
	if !exists {
		return nil, errors.New("jwt claims not found in context")
	}

	// check if the claims are of type jwt.MapClaims
	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid jwt claims format")
	}

	return claims, nil
}

// Get a specific claim from the JWT claims in the context
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

// getUsername retrieves the username from the JWT claims in the context
func GetUsername(c *gin.Context) (string, error) {
	return getClaim(c, "sub")
}

// retrieve the user role from the JWT claims in the context
func GetRole(c *gin.Context) (string, error) {
	return getClaim(c, "role")
}

// retrieve the device ID from the JWT claims in the context
func GetDeviceID(c *gin.Context) (string, error) {
	return getClaim(c, "device_id")
}

// retrieve the CSRF token from the JWT claims in the context
func GetCSRFToken(c *gin.Context) (string, error) {
	return getClaim(c, "csrf_token")
}

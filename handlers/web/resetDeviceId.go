package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

// ResetDeviceRequest defines the structure for the reset device ID request
type ResetDeviceRequest struct {
	Username string `json:"username" binding:"required"`
}

// function ResetDeviceID handles the request to reset the device ID for a user
func ResetDeviceID (c *gin.Context)  {

	// Validate the request payload
	var req ResetDeviceRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if the user exists in the database
	// If the user does not exist, return an error
	// If the user exists, proceed to reset the device ID
	if err := config.DB.Model(&models.User{}).Where("username = ?", req.Username).Update("device_id", nil).Update("device_id_updated_at", nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset device ID"})
		return
	}

	// If the update is successful, return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Device ID reset successfully"})
}

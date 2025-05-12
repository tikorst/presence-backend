package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/models"
)

type ResetDeviceRequest struct {
	Username string `json:"username" binding:"required"`
}

func ResetDeviceID(c *gin.Context) {
	var req ResetDeviceRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := config.DB.Model(&models.User{}).Where("username = ?", req.Username).Update("device_id", nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset device ID"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Device ID reset successfully"})
}

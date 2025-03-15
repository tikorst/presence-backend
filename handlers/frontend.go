package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tikorst/siatma-backend/config"
	"github.com/tikorst/siatma-backend/models"
)

func Frontend() gin.HandlerFunc {
	return func(c *gin.Context) {
		var classes []models.Kelas
		config.DB.Find(&classes)
		c.HTML(http.StatusOK, "index.html", gin.H{"classes": classes})
	}



}

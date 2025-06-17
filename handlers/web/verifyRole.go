package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/presence-backend/helpers"
)

func VerifyRole(c *gin.Context) {
	role, _ := helpers.GetRole(c)

	c.JSON(http.StatusOK, gin.H{"role": role})
}

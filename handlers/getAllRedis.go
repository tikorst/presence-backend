package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tikorst/siatma-backend/config"
)

func GetAll(c *gin.Context) {
	// Get all data from Redis
	fmt.Println("Getting all data from Redis")
	val, err := config.RedisDB.Keys(config.Ctx, "*").Result()
	if err != nil {
		fmt.Println("Error getting data from Redis:", err)
	}
	fmt.Println("Data from Redis:", val)
	c.JSON(200, gin.H{"data": val})
}

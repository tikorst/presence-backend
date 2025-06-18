package web

import "github.com/gin-gonic/gin"

func Logout(c *gin.Context) {
	c.Header("Set-Cookie", "token=; Path=/; Domain=.tikorst.cloud; Max-Age=-1; HttpOnly; SameSite=Lax; Secure")
	c.JSON(200, gin.H{"message": "Logout successful"})
}

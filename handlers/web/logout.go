package web

import "github.com/gin-gonic/gin"

// Function to handle user logout
func Logout(c *gin.Context) {

	// Clear the token cookie by setting it to an empty value and a past expiration date
	c.Header("Set-Cookie", "token=; Path=/; Domain=.tikorst.cloud; Max-Age=-1; HttpOnly; SameSite=Lax; Secure")
	c.JSON(200, gin.H{"message": "Logout successful"})
}

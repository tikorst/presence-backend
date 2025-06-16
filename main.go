package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tikorst/presence-backend/config"
	"github.com/tikorst/presence-backend/handlers/mobile"
	"github.com/tikorst/presence-backend/handlers/web"
	"github.com/tikorst/presence-backend/middleware"
)

func init() {
	godotenv.Load()
	config.InitTLS()
	config.ConnectDB()
	config.ConnectRedis()
	config.ConnectStorage()
}
func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://presence-web.tikorst.cloud"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API
	r.POST("/login", mobile.Login)

	// Web
	r.POST("/web/login", web.Login)
	protectedWeb := r.Group("/web")
	{
		protectedWeb.Use(middleware.Web())
		protectedWeb.GET("/classes", web.GetClasses)
		protectedWeb.GET("/classes/:classID/meetings", web.GetMeetings)
		protectedWeb.GET("/generate_qr/:classID/:meetingID", web.GenerateQR)
		protectedWeb.GET("/attendance/:classID/:meetingID", web.GetPresenceData)
		protectedWeb.GET("verify-role", web.VerifyRole)
		protectedWeb.POST("/attendance/:classID/:meetingID", web.ManualAttendance)

		adminRoutes := protectedWeb.Group("/admin").Use(middleware.AdminOnly())
		{
			adminRoutes.GET("/users", web.GetUsers)
			adminRoutes.POST("/reset", web.ResetDeviceID)
		}
	}

	// Protected routes
	protected := r.Group("/api").Use(middleware.Auth())
	{
		protected.POST("/presence", mobile.ValidateQr)
		protected.GET("/validate", mobile.ValidateToken)
		protected.GET("/class", mobile.GetSchedules)
		protected.GET("/attendance", mobile.GetAttendance)
		protected.GET("/semester", mobile.GetSemester)
		protected.GET("/grade", mobile.GetGrade)
		protected.GET("/allgrade", mobile.GetAllGrade)
		protected.GET("/profile", mobile.GetProfile)
	}
	port := os.Getenv("PORT")
	certPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("KEY_PATH")

	r.RunTLS("0.0.0.0:"+port, certPath, keyPath)
}

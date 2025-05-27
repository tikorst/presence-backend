package main

import (
	"fmt"
	"os"
	"runtime"
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
	godotenv.Load()
}
func main() {
	runtime.GOMAXPROCS(2)
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// Public routes
	if config.RedisDB == nil {
		fmt.Println("RedisDB is nil! Initialization failed.")
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://presence-web.tikorst.cloud"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API
	r.POST("/login", mobile.Login())

	// Web
	r.POST("/web/login", web.Login())
	protectedWeb := r.Group("/web")
	{
		protectedWeb.Use(middleware.Web())
		protectedWeb.GET("/classes", web.Classes())
		protectedWeb.GET("/classes/:classID/meetings", web.Meetings())
		protectedWeb.GET("/generate_qr/:classID/:meetingID", web.GenerateQR)
		protectedWeb.GET("/attendance/:classID/:meetingID", web.PresenceData)
		protectedWeb.GET("verify-role", web.VerifyRole())

		adminRoutes := protectedWeb.Group("/admin").Use(middleware.AdminOnly())
		{
			adminRoutes.GET("/users", web.GetUsers())
			adminRoutes.POST("/reset", web.ResetDeviceID)
		}
	}

	// Protected routes
	protected := r.Group("/api").Use(middleware.Auth())
	{
		protected.POST("/presence", mobile.ValidateQr)
		protected.GET("/validate", mobile.ValidateToken())
		protected.GET("/kelas", mobile.Jadwal)
		protected.GET("/attendance", mobile.Attendance)
		protected.GET("/semester", mobile.Semester)
		protected.GET("/grade", mobile.Grade)
		protected.GET("/allgrade", mobile.AllGrade)
		protected.GET("/profile", mobile.Profile)
	}
	port := os.Getenv("PORT")
	r.Run("0.0.0.0:" + port) // Jalankan di port 8080
}

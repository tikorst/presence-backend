// package main

// import (
// 	"fmt"
// 	"runtime"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/joho/godotenv"
// 	"github.com/tikorst/presence-backend/config"
// 	"github.com/tikorst/presence-backend/handlers"
// )

// func init() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		panic("Failed to load .env file")
// 	}
// 	config.InitTLS()
// 	config.ConnectDB()
// 	fmt.Println("🖥️ CPU Cores Used:", runtime.NumCPU())
// 	fmt.Println("🔄 Goroutines at startup:", runtime.NumGoroutine())
// }

// func main() {
// 	runtime.GOMAXPROCS(20)
// 	r := fiber.New()
// 	// r(gin.ReleaseMode)
// 	// Public routes
// 	r.Get("/login", handlers.Login1())
// 	// Protected routes
// 	// protected := r.Group("/api").Use(middleware.Auth())
// 	// {
// 	// 	protected.POST("/presensi", handlers.ScanQR())
// 	// }
// 	go func() {
// 		for {
// 			fmt.Println("🔄 Goroutines:", runtime.NumGoroutine())
// 			fmt.Println("🖥️ CPU Cores Used:", runtime.GOMAXPROCS(0))
// 			time.Sleep(1000 * time.Millisecond) // Log every 5 seconds
// 		}
// 	}()
// 	r.Listen(":8080") // Jalankan di port 8080
// }

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
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API
	// r.GET("/ping", handlers.Ping())
	r.POST("/login", mobile.Login())
	// r.GET("/login", handlers.Login2())
	// r.GET("/frontend", handlers.Frontend())
	// r.GET("/getAllRedis", handlers.GetAll)

	// Web
	r.POST("/web/login", web.Login())
	protectedWeb := r.Group("/web").Use(middleware.Web())
	{
		protectedWeb.GET("/classes", web.Classes())
		protectedWeb.GET("/classes/:classID/meetings", web.Meetings())
		protectedWeb.GET("/generate_qr/:classID/:meetingID", web.GenerateQR)
		protectedWeb.GET("/attendance/:classID/:meetingID", web.PresenceData)
	}

	// Protected routes
	protected := r.Group("/api").Use(middleware.Auth())
	{
		protected.POST("/presence", mobile.ValidateQr)
		protected.GET("/validate", mobile.ValidateToken())
		protected.GET("/kelas", mobile.Jadwal)
		protected.GET("/attendance", mobile.Attendance)
	}
	port := os.Getenv("PORT")
	r.Run("0.0.0.0:" + port) // Jalankan di port 8080
}

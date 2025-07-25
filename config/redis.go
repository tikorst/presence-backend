package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx     = context.Background()
	RedisDB *redis.Client
)

func ConnectRedis() {

	// Load Redis connection details from environment variables
	redis_addr := os.Getenv("REDIS_ADDR")
	redis_pass := os.Getenv("REDIS_PASSWORD")

	// Connect to Redis
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_pass,
		DB:       0,
	})

	// Check if RedisDB is nil
	if RedisDB == nil {
		fmt.Println("RedisDB is nil")
	}

	// Check if RedisDB is not nil
	if RedisDB != nil {
		fmt.Println("RedisDB is not nil")
	}
	// Test the connection

	// Use Ping to check the connection
	result, err := RedisDB.Ping(Ctx).Result()
	RedisDB.Set(Ctx, "qr", 1, 150*time.Second)


	// If there is an error, log it and exit
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}

	// If the connection is successful, print the result
	fmt.Println("Connected to Redis! with result:", result)

}

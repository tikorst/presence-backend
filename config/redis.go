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
	redis_addr := os.Getenv("REDIS_ADDR")
	redis_pass := os.Getenv("REDIS_PASSWORD")
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_pass,
		DB:       0,
	})
	if RedisDB == nil {
		fmt.Println("RedisDB is nil")
	}
	if RedisDB != nil {
		fmt.Println("RedisDB is not nil")
	}
	// Test the connection
	result, err := RedisDB.Ping(Ctx).Result()
	RedisDB.Set(Ctx, "qr", 1, 150*time.Second)

	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	fmt.Println("Connected to Redis! with result:", result)

}

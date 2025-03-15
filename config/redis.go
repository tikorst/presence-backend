package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx     = context.Background()
	RedisDB *redis.Client
)

func ConnectRedis() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "tiko07",
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

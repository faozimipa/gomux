package redismanager

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v7"
)

func NewRedisDB(host, port, password string) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	return redisClient
}

func InitRedisClient() redis.Client {
	// client := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "",
	// 	DB:       0, //default
	// })

	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASS")

	redisClient := NewRedisDB(redis_host, redis_port, redis_password)

	pong, err := redisClient.Ping().Result()
	if err != nil {
		fmt.Println("Cannot Initialize Redis Client ", err)
	}
	fmt.Println("Redis Client Successfully Initialized . . .", pong)

	return *redisClient
}

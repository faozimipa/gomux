package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/faozimipa/gomux/api/responses"
	"github.com/faozimipa/gomux/api/utils/redismanager"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func (server *Server) TestRedis(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Sampe sini kok.")
	c := redismanager.InitRedisClient()
	defer c.Close()

	err := c.Set(ctx, "key", "value", 0).Err()

	if err != nil {
		fmt.Println(err)
	}

	val, err := c.Do(ctx, "get", "key").Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("key does not exists")
			return
		}
		fmt.Println(err)
	}

	responses.JSON(w, http.StatusOK, val)

}

func (server *Server) SetData(w http.ResponseWriter, r *http.Request) {

	c := redismanager.InitRedisClient()
	defer c.Close()
	// SET key value EX 10 NX
	set, err := c.SetNX(ctx, "key1", "value ini muncul 10s", 10*time.Second).Result()

	if err != nil {
		if err != redis.Nil {
			fmt.Println("key does not exists")
			return
		}
		fmt.Println(err)
	}
	responses.JSON(w, http.StatusOK, set)
}

func (server *Server) GetData(w http.ResponseWriter, r *http.Request) {

	c := redismanager.InitRedisClient()
	defer c.Close()

	val, err := c.Do(ctx, "get", "key1").Result()

	if err != nil {
		if err != redis.Nil {
			fmt.Println("key does not exists")
			return
		}
		fmt.Println(err)
	}

	responses.JSON(w, http.StatusOK, val)
}

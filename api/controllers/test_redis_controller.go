package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/faozimipa/gomux/api/responses"
	"github.com/faozimipa/gomux/api/utils/rabbitmq"
	"github.com/faozimipa/gomux/api/utils/redismanager"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
)

var ctx = context.Background()

func (server *Server) TestRedis(w http.ResponseWriter, r *http.Request) {
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

	//test mq
	cr, err := rabbitmq.ConnectMQ()

	if err != nil {
		log.Fatalf("conn mq error  ")
	}
	defer cr.Close()

	ch, err := cr.Channel()
	if err != nil {
		log.Fatalf("can not create channel")
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"TestQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	// We can print out the status of our Queue here
	// this will information like the amount of messages on
	// the queue
	fmt.Println(q)
	// Handle any errors if we were unable to create the queue
	if err != nil {
		fmt.Println(err)
	}

	// attempt to publish a message to the queue!
	err = ch.Publish(
		"",
		"TestQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello World"),
		},
	)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Published Message to Queue")

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

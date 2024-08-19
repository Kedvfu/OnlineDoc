package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
)

var RedisClient *redis.Client

func InitializeRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := RedisClient.Ping().Result()
	if err != nil {
		log.Printf("Unable to connect to redis: %s", err)
		panic(err)
		return
	}
	log.Println(pong, err)

}
func Ping(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "pong",
	})
}

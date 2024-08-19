package handlers

import (
	"fmt"
	"github.com/go-redis/redis"
	"testing"
)

func TestRedis(t *testing.T) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	RedisClient.Set("key", "value", 0)
	val, err := RedisClient.Get("key").Result()
	if err != nil {
		t.Error(err)
	}
	if val != "value" {
		t.Error("Redis set and get failed")
	}
	fmt.Println(val)
}

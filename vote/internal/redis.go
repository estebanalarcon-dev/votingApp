package internal

import (
	"fmt"
	"github.com/go-redis/redis"
)

type RedisClient interface {
	RPush(key string, values ...interface{}) error
}

type redisClient struct {
	client *redis.Client
}

func NewRedisClient() RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return redisClient{client: client}
}

func (r redisClient) RPush(key string, values ...interface{}) error {
	res := r.client.RPush(key, values...)
	return res.Err()
}

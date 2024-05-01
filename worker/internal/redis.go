package internal

import (
	"fmt"
	"github.com/go-redis/redis"
)

type RedisClient interface {
	RPop(key string) (string, error)
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

func (r redisClient) RPop(key string) (string, error) {
	res := r.client.RPop(key)
	return res.Result()
}

package redisDB

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("Redis bilan ulanishda xato: %v", err)
		return nil, err
	}
	log.Printf("Redis serveriga muvaffaqiyatli ulandi: %s", pong)

	return rdb, nil
}

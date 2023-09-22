package database

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisInstance []*redis.Client

func InitRedis() error {
	for i := 0; i < 1; i++ {
		addr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
		client := redis.NewClient(
			&redis.Options{
				Addr:     addr,
				Password: os.Getenv("REDIS_PASS"),
				DB:       i,
			},
		)

		// validate connection
		_, err := client.Ping(context.Background()).Result()
		if err != nil {
			return err
		}

		RedisInstance = append(RedisInstance, client)
	}
	return nil
}

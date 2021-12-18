package database

import (
	"github.com/go-redis/redis/v7"
)

func Redis() *redis.Client {
	// Redis Server TEST at Ping
	dsn := ""
	client := redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	return client
}

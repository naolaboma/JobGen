package infrastructure

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

// NewRedisClient initializes a Redis v8 client with the given address.
// addr format is typically "host:port". Uses DB 0 by default.
func NewRedisClient(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Printf("warning: unable to ping redis at %s: %v", addr, err)
	}
	return client
}

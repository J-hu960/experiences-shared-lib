package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}

func PublishEvent(client *redis.Client, channel string, message string) error {
	pubsub := client.Publish(context.Background(), channel, message)
	if pubsub.Err() != nil {
		return fmt.Errorf("error publishing message: %v", pubsub.Err())
	}

	return nil
}

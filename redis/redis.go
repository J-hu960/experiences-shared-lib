package redis_internal

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type EventHandler interface {
	HandleEvent(msg *redis.Message)
}

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}

func PublishEvent(client *redis.Client, channel string, message interface{}) error {
	pubsub := client.Publish(context.Background(), channel, message)
	if pubsub.Err() != nil {
		return fmt.Errorf("error publishing message: %v", pubsub.Err())
	}

	return nil
}

func SubscribeToEvents(client *redis.Client, channel string, handler EventHandler) {
	ctx := context.Background()
	pubsub := client.Subscribe(ctx, channel) //aqui deberiamos pasarle un array de channels a los que se va a suscribir?
	defer pubsub.Close()

	// Bloquea y espera mensajes
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			fmt.Printf("Error receiving message: %v\n", err)
			continue
		}
		handler.HandleEvent(msg)

		fmt.Printf("Message %s recieved ", msg.Channel)
	}
}

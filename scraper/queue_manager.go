package scraper

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

type QueueManager struct {
	client *redis.Client
}

func NewQueueManager(redisURL string) *QueueManager {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	return &QueueManager{client: rdb}
}

func (qm *QueueManager) Push(ctx context.Context, queue string, value string) {
	err := qm.client.LPush(ctx, queue, value).Err()
	if err != nil {
		log.Println("Error pushing to queue:", err)
	}
}

func (qm *QueueManager) Pop(ctx context.Context, queue string) string {
	val, err := qm.client.RPop(ctx, queue).Result()
	if err != nil {
		return ""
	}
	return val
}

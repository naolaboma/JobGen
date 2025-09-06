package infrastructure

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type QueueService interface {
	Enqueue(jobID string) error
	Dequeue() (string, error)
}

type redisQueueService struct {
	client    *redis.Client
	queueName string
}

func NewQueueService(client *redis.Client, queueName string) QueueService {
	return &redisQueueService{
		client:    client,
		queueName: queueName,
	}
}

func (q *redisQueueService) Enqueue(jobID string) error {
	return q.client.LPush(context.Background(), q.queueName, jobID).Err()
}

func (q *redisQueueService) Dequeue() (string, error) {
	// Blocking pop: waits forever until a job is available.
	result, err := q.client.BRPop(context.Background(), 0*time.Second, q.queueName).Result()
	if err != nil {
		return "", err
	}
	return result[1], nil // result is a slice [queueName, value]
}

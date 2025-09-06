package infrastructure

import (
	"errors"
)

// inMemoryQueueService is a simple channel-backed queue used as a fallback
// when Redis is not available. It satisfies the QueueService interface.
type inMemoryQueueService struct {
	ch chan string
}

// NewInMemoryQueueService creates a new in-memory queue with a reasonable buffer.
func NewInMemoryQueueService(buffer int) QueueService {
	if buffer <= 0 {
		buffer = 100
	}
	return &inMemoryQueueService{ch: make(chan string, buffer)}
}

func (q *inMemoryQueueService) Enqueue(jobID string) error {
	select {
	case q.ch <- jobID:
		return nil
	default:
		// If buffer is full, try a blocking send to ensure delivery.
		q.ch <- jobID
		return nil
	}
}

func (q *inMemoryQueueService) Dequeue() (string, error) {
	jobID, ok := <-q.ch
	if !ok {
		return "", errors.New("queue closed")
	}
	return jobID, nil
}

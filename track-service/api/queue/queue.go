package queue

import (
	"context"
	"fmt"
	"tracker/internal/usecase/queueservice"
)

type Queue struct {
	*queueservice.QueueService
}

func New(srv *queueservice.QueueService) *Queue {
	return &Queue{srv}
}

func (q *Queue) WriteMessages(ctx context.Context, msg []byte) error {
	err := q.Write(ctx, msg)
	if err != nil {
		return fmt.Errorf("queue: %w", err)
	}
	return nil
}

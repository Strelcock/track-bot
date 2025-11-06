package queueservice

import (
	"context"
	"fmt"
)

type QueueService struct {
	IQueueRepo
}

func New(repo IQueueRepo) *QueueService {
	return &QueueService{repo}
}

func (q *QueueService) Write(ctx context.Context, msg []byte) error {
	err := q.IQueueRepo.Write(ctx, msg)
	if err != nil {
		return fmt.Errorf("service: %w", err)
	}
	return nil
}

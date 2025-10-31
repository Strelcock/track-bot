package queue

import "context"

type IQueueService interface {
	Read(ctx context.Context) ([]byte, error)
}

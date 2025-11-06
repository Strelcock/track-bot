package queueservice

import "context"

type IQueueRepo interface {
	Write(ctx context.Context, msg []byte) error
}

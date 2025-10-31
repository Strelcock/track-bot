package kafkaservice

import "context"

type IKafkaRepo interface {
	Read(ctx context.Context) ([]byte, error)
}

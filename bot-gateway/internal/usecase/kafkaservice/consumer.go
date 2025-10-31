package kafkaservice

import (
	"context"
)

type KafkaService struct {
	IKafkaRepo
}

func New(repo IKafkaRepo) *KafkaService {
	return &KafkaService{repo}
}

func (s *KafkaService) Read(ctx context.Context) ([]byte, error) {
	bytes, err := s.IKafkaRepo.Read(ctx)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

package kservice

import (
	"context"
	"core-service/internal/domain/track"
	"log"
)

type ProducerService struct {
	track.TrackProducer
}

func NewProd(producer track.TrackProducer) *ProducerService {
	return &ProducerService{producer}
}

func (p *ProducerService) SendMessage(ctx context.Context, msg []byte) error {
	err := p.Write(ctx, msg)
	if err != nil {
		return err
	}
	log.Println("Message sent")
	return nil
}

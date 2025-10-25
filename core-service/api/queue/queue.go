package queue

import (
	"context"
	"core-service/internal/usecase/kservice"
	"core-service/pkg/json"
	"log"
	"time"
)

type Queue struct {
	*kservice.ConsumerService
	*kservice.ProducerService
}

func New(cons *kservice.ConsumerService, prod *kservice.ProducerService) *Queue {
	return &Queue{cons, prod}
}

func (q *Queue) ServeMessages() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	to, err := q.ConsumerService.Read(ctx)
	if err != nil {
		log.Println(err)
	}

	bytes, err := json.Marshal(to)
	if err != nil {
		log.Println(err)
	}

	go q.ProducerService.SendMessage(ctx, bytes)
}

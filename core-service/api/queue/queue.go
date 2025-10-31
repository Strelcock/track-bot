package queue

import (
	"context"
	"core-service/internal/usecase/kservice"
	"core-service/pkg/json"
	"log"
)

type Queue struct {
	*kservice.ConsumerService
	*kservice.ProducerService
}

func New(cons *kservice.ConsumerService, prod *kservice.ProducerService) *Queue {
	return &Queue{cons, prod}
}

func (q *Queue) ServeMessages() {

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	for {
		func() {
			defer func() { // если приходит сообщение о посылке которой
				recover() // у пользователя нет
			}() //притянуто за уши для разработки
			to, err := q.ConsumerService.Read(context.Background())
			if err != nil {
				log.Println(err)
			}

			if to.Status.ID == "0" {
				return
			}

			bytes, err := json.Marshal(to)
			if err != nil {
				log.Println(err)
			}

			go q.ProducerService.SendMessage(context.Background(), bytes)
		}()
	}
}

package queue

import (
	"bot-gateway/api/bot"
	"context"
	"encoding/json"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Queue struct {
	IQueueService
}

func New(srv IQueueService) *Queue {
	return &Queue{srv}
}

func (q *Queue) Read(b *bot.Bot) {
	ctx := context.Background()
	for {
		bytes, err := q.IQueueService.Read(ctx)

		if err != nil {
			log.Println(err)
			return
		}

		var status StatusTo
		err = json.Unmarshal(bytes, &status)

		if err != nil {
			log.Println(err)
			return
		}

		statusStr := fmt.Sprintf("Заказ номер %s %s\nДата: %v", status.Status.ID, status.Status.Status, status.Status.UpdatedAt)

		msg := tgbotapi.NewMessage(status.To, statusStr)
		_, err = b.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

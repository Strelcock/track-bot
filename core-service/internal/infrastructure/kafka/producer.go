package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	*kafka.Writer
}

func NewProd(brokers []string, topic string) *Producer {
	return &Producer{kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})}
}

func (p *Producer) Write(ctx context.Context, msg []byte) error {
	err := p.WriteMessages(ctx, kafka.Message{
		Value: msg,
	})

	if err != nil {
		return err
	}
	return nil
}

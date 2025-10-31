package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	*kafka.Reader
}

func New(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (c *Consumer) Read(ctx context.Context) ([]byte, error) {
	msg, err := c.Reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}

	return msg.Value, nil
}

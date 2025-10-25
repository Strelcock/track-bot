package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	*kafka.Reader
}

func NewCons(brokers []string, topic string, groupID string) *Consumer {
	return &Consumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (c *Consumer) Read(ctx context.Context) ([]byte, error) {
	msg, err := c.ReadMessage(ctx)
	if err != nil {
		return []byte{}, err
	}

	return msg.Value, nil

}

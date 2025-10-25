package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	*kafka.Reader
}

func New(config kafka.ReaderConfig) *KafkaConsumer {
	return &KafkaConsumer{
		Reader: kafka.NewReader(config),
	}
}

func (c *KafkaConsumer) Read(ctx context.Context) ([]byte, error) {
	msg, err := c.ReadMessage(ctx)
	if err != nil {
		return []byte{}, err
	}

	return msg.Value, nil

}

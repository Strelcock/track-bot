package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	*kafka.Reader
}

func New(brokers []string, topic, groupID string) *Consumer {
	time.Sleep(10 * time.Second)
	conn, err := kafka.DialLeader(context.Background(), "tcp", brokers[0], topic, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	controller, _ := conn.Controller()
	log.Print(controller.Host, controller.Port)

	return &Consumer{
		kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       topic,
			GroupID:     groupID,
			StartOffset: kafka.FirstOffset,
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

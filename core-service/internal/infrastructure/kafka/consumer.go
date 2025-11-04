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

func NewCons(brokers []string, topic string, groupID string) *Consumer {
	time.Sleep(10 * time.Second)
	conn, err := kafka.DialLeader(context.Background(), "tcp", brokers[0], topic, 0)
	if err != nil {
		log.Printf("failed to dial leader: %s\n", err.Error())
	}
	contoller, _ := conn.Controller()
	log.Println(contoller.Host, contoller.Port)
	defer conn.Close()

	return &Consumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       topic,
			GroupID:     groupID,
			StartOffset: kafka.FirstOffset,
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

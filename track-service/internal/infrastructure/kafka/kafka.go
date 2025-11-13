package kafka

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	*kafka.Writer
}

func New(brokers []string, topic string) *Producer {
	for i := range 5 {
		time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
		conn, err := kafka.DialLeader(context.Background(), "tcp", brokers[0], topic, 0)
		if err == nil {
			controller, _ := conn.Controller()
			log.Print(controller.Host, controller.Port)
			conn.Close()
			break
		} else if i == 4 {
			log.Fatal(err)
		}
	}
	return &Producer{
		kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic,
		}),
	}
}

func (p *Producer) Write(ctx context.Context, msg []byte) error {
	err := p.WriteMessages(ctx, kafka.Message{
		Value: msg,
	})
	if err != nil {
		return fmt.Errorf("producer: write: %w", err)
	}
	return nil
}

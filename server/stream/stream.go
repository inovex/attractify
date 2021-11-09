package stream

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Stream struct {
	brokers []string
	w       *kafka.Writer
}

func New(brokers []string, topic string) *Stream {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Async:    true,
		Balancer: &kafka.RoundRobin{},
	})

	return &Stream{
		brokers: brokers,
		w:       w,
	}
}

func (s Stream) Write(ctx context.Context, key, message []byte) error {
	return s.w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   key,
			Value: message,
		},
	)
}

func (s Stream) Close() error {
	return s.w.Close()
}

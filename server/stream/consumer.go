package stream

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

func (s Stream) Consume(ctx context.Context, topic, groupID string, cb func(*Msg) error) error {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   s.brokers,
		Topic:     topic,
		Partition: 0,
		MinBytes:  1,
		MaxBytes:  10e6,
		GroupID:   groupID,
		// RebalanceTimeout:       time.Second,
		WatchPartitionChanges:  true,
		PartitionWatchInterval: time.Second,
	})
	defer r.Close()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return nil
		default:
			msg, _ := NewMsg(m)
			if err := cb(msg); err != nil {
				continue
			}

			r.CommitMessages(ctx, m)
		}
	}
}

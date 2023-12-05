package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/somatom98/badges/domain"
)

type EventKafkaConsumer struct {
}

func NewEventKafkaConsumer() *EventKafkaConsumer {
	return &EventKafkaConsumer{}
}

func (c *EventKafkaConsumer) Consume(ctx context.Context) (<-chan *domain.Event, error) {
	ch := make(chan *domain.Event)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "badge-events",
		Partition: 0,
		MaxBytes:  10e6,
	})

	r.SetOffset(-1)

	go func() {
		defer r.Close()
		for {
			m, err := r.FetchMessage(ctx)
			if err != nil {
				break
			}
			fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

			var event domain.Event
			err = json.Unmarshal(m.Value, &event)
			if err != nil {
				log.Error().
					Err(err).
					Msg("kafka deserialization error")
				continue
			}

			ch <- &event

			r.CommitMessages(ctx, m)
		}
	}()

	return ch, nil
}

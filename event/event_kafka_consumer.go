package event

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/somatom98/badges/config"
	"github.com/somatom98/badges/domain"
)

type EventKafkaConsumer struct {
	options config.KafkaOptions
}

func NewEventKafkaConsumer(options config.KafkaOptions) *EventKafkaConsumer {
	return &EventKafkaConsumer{
		options: options,
	}
}

func (c *EventKafkaConsumer) Consume(ctx context.Context, handler *func(context.Context, domain.Event) error) (<-chan *domain.Event, error) {
	ch := make(chan *domain.Event)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   c.options.Brokers,
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

			var event domain.Event
			err = json.Unmarshal(m.Value, &event)
			if err != nil {
				log.Error().
					Err(err).
					Msg("kafka deserialization error")
				r.CommitMessages(ctx, m)
				continue
			}

			if handler != nil {
				function := *handler
				err = function(ctx, event)
				if err != nil {
					log.Error().
						Err(err).
						Msg("handler error")
					continue
				}
			}

			ch <- &event
			r.CommitMessages(ctx, m)
		}
	}()

	return ch, nil
}

package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/somatom98/badges/domain"
)

func main() {
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "badge-events",
		Balancer: &kafka.LeastBytes{},
	}

	event := domain.Event{
		ID:   "asd",
		UID:  "user2",
		Type: domain.EventTypeIn,
		Date: time.Now().UTC(),
	}
	marshalledEvent, err := json.Marshal(event)
	if err != nil {
		log.Print("Panic")
		panic(err)
	}

	log.Print(event.Date.String())
	err = w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("id1"),
			Value: marshalledEvent,
		},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
	log.Print("Sent")

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

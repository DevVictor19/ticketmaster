package cdc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type EventDbConsumer struct {
	bootstrapServers string
	shutdown         chan uint8
}

func NewEventDbConsumer(bootstrapServers string) *EventDbConsumer {
	return &EventDbConsumer{
		bootstrapServers: bootstrapServers,
		shutdown:         make(chan uint8),
	}
}

var (
	eventTopic = "pg.public.events"
	venueTopic = "pg.public.venues"
)

func (ec *EventDbConsumer) Start() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": ec.bootstrapServers,
		"group.id":          "search-api",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}

	err = c.SubscribeTopics([]string{
		eventTopic,
		venueTopic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe consumer: %s", err)
	}

	go func() {
		run := true
		for run {
			select {
			case <-ec.shutdown:
				run = false
				c.Close()
				slog.Info("CDC consumer shutdown complete")
				ec.shutdown <- 1
				return

			default:
				msg, err := c.ReadMessage(time.Second)
				if err == nil {
					ec.processMessage(msg)
				} else if !err.(kafka.Error).IsTimeout() {
					slog.Error("Consumer error", "error", err, "message", msg)
				}
			}
		}
	}()
}

func (ec *EventDbConsumer) Stop(ctx context.Context) {
	ec.shutdown <- 0
	for {
		select {
		case <-ec.shutdown:
			return
		case <-ctx.Done():
			slog.Warn("Timeout waiting for CDC consumer to shutdown")
			return
		}
	}
}

func (ec *EventDbConsumer) processMessage(msg *kafka.Message) {
	if msg.Value == nil {
		slog.Debug("Tombstone event received, skipping", "topic", *msg.TopicPartition.Topic)
		return
	}

	topic := *msg.TopicPartition.Topic

	switch topic {
	case eventTopic:
		ec.handleEventPayload(msg.Value)
	case venueTopic:
		ec.handleVenuePayload(msg.Value)
	default:
		slog.Warn("Unknown topic", "topic", topic)
	}
}

func (ec *EventDbConsumer) handleEventPayload(value []byte) {
	envelope, err := unmarshalEnvelope[Event](value)
	if err != nil {
		slog.Error("error unmarshaling event payload", "error", err)
		return
	}

	switch envelope.Payload.Op {
	case CreateOp:
		fmt.Printf("New event created: %+v\n", envelope.Payload.After)
	case UpdateOp:
		fmt.Printf("Event updated from %+v to %+v\n", envelope.Payload.Before, envelope.Payload.After)
	case DeleteOp:
		fmt.Printf("Event deleted: %+v\n", envelope.Payload.Before)
	}

}

func (ec *EventDbConsumer) handleVenuePayload(value []byte) {
	envelope, err := unmarshalEnvelope[Venue](value)
	if err != nil {
		slog.Error("error unmarshaling venue payload", "error", err)
		return
	}

	switch envelope.Payload.Op {
	case CreateOp:
		fmt.Printf("New venue created: %+v\n", envelope.Payload.After)
	case UpdateOp:
		fmt.Printf("Venue updated from %+v to %+v\n", envelope.Payload.Before, envelope.Payload.After)
	case DeleteOp:
		fmt.Printf("Venue deleted: %+v\n", envelope.Payload.Before)
	}
}

func unmarshalEnvelope[T any](value []byte) (*DebeziumEnvelope[T], error) {
	var envelope DebeziumEnvelope[T]
	err := json.Unmarshal(value, &envelope)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling message: %w", err)
	}
	return &envelope, nil
}

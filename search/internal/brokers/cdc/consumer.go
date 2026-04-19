package cdc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/DevVictor19/search/internal/services"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type EventDbConsumer struct {
	bootstrapServers string
	shutdown         chan uint8
	searchEngineSvc  services.SearchEngineService
}

func NewEventDbConsumer(
	bootstrapServers string,
	searchEngineSvc services.SearchEngineService) *EventDbConsumer {

	return &EventDbConsumer{
		bootstrapServers: bootstrapServers,
		shutdown:         make(chan uint8),
		searchEngineSvc:  searchEngineSvc,
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
		ec.handleEventMsg(msg.Value)
	case venueTopic:
		ec.handleVenueMsg(msg.Value)
	default:
		slog.Warn("Unknown topic", "topic", topic)
	}
}

func (ec *EventDbConsumer) handleEventMsg(value []byte) {
	envelope, err := unmarshalEnvelope[Event](value)
	if err != nil {
		slog.Error("error unmarshaling event payload", "error", err)
		return
	}

	operation := envelope.Payload.Op
	if operation == CreateOp || operation == UpdateOp {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = ec.searchEngineSvc.UpsertEventIdx(ctx, &services.EventDoc{
			ID:          envelope.Payload.After.ID,
			Name:        envelope.Payload.After.Name,
			Description: envelope.Payload.After.Description,
			Date:        envelope.Payload.After.Date,
			VenueID:     envelope.Payload.After.VenueID,
		})
		if err != nil {
			slog.Error(
				"error creating event index", "error", err, "event_id",
				envelope.Payload.After.ID,
			)
		}
		return
	}

	if operation == DeleteOp {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = ec.searchEngineSvc.DeleteEventIdx(ctx, envelope.Payload.Before.ID)
		if err != nil {
			slog.Error("error deleting event index", "error", err, "event_id", envelope.Payload.Before.ID)
		}
		return
	}
}

func (ec *EventDbConsumer) handleVenueMsg(value []byte) {
	envelope, err := unmarshalEnvelope[Venue](value)
	if err != nil {
		slog.Error("error unmarshaling venue payload", "error", err)
		return
	}

	operation := envelope.Payload.Op
	if operation == UpdateOp {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = ec.searchEngineSvc.UpdateEventLocationIdx(
			ctx,
			envelope.Payload.After.ID,
			envelope.Payload.After.Location,
		)
		if err != nil {
			slog.Error(
				"error updating venue index", "error", err, "venue_id",
				envelope.Payload.After.ID,
			)
		}
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

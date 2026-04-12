package cdc

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type EventManagementDbConsumer struct {
	BootstrapServers string
}

func NewEventManagementDbConsumer(bootstrapServers string) *EventManagementDbConsumer {
	return &EventManagementDbConsumer{
		BootstrapServers: bootstrapServers,
	}
}

var (
	ticketTopic         = "pg.public.tickets"
	eventTopic          = "pg.public.events"
	venueTopic          = "pg.public.venues"
	performerTopic      = "pg.public.performers"
	eventPerformerTopic = "pg.public.event_performers"
)

func (ec *EventManagementDbConsumer) Start() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": ec.BootstrapServers,
		"group.id":          "search-api",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}

	err = c.SubscribeTopics([]string{
		eventPerformerTopic,
		eventTopic,
		performerTopic,
		ticketTopic,
		venueTopic}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe consumer: %s", err)
	}

	run := true

	for run {
		msg, err := c.ReadMessage(time.Second)
		if err == nil {
			processMessage(msg)
		} else if !err.(kafka.Error).IsTimeout() {
			slog.Error("Consumer error", "error", err, "message", msg)
		}
	}

	c.Close()
}

func processMessage(msg *kafka.Message) {
	if msg.Value == nil {
		slog.Debug("Tombstone event received, skipping", "topic", *msg.TopicPartition.Topic)
		return
	}

	topic := *msg.TopicPartition.Topic

	switch topic {
	case ticketTopic:
		handleTicketEvent(msg.Value)
	case eventTopic:
		handleEventEvent(msg.Value)
	case venueTopic:
		handleVenueEvent(msg.Value)
	case performerTopic:
		handlePerformerEvent(msg.Value)
	case eventPerformerTopic:
		handleEventPerformerEvent(msg.Value)
	default:
		slog.Warn("Unknown topic", "topic", topic)
	}
}

func handleTicketEvent(value []byte) {
	envelope, err := unmarshalEnvelope[Ticket](value)
	if err != nil {
		slog.Error("error unmarshaling ticket event", "error", err)
		return
	}

	switch envelope.Payload.Op {
	case "c": // create
		fmt.Printf("New ticket created: %+v\n", envelope.Payload.After)
	case "u": // update
		fmt.Printf("Ticket updated from %+v to %+v\n", envelope.Payload.Before, envelope.Payload.After)
	case "d": // delete
		fmt.Printf("Ticket deleted: %+v\n", envelope.Payload.Before)
	}
}

func handleEventEvent(value []byte) {}

func handleVenueEvent(value []byte) {}

func handlePerformerEvent(value []byte) {}

func handleEventPerformerEvent(value []byte) {}

func unmarshalEnvelope[T any](value []byte) (*DebeziumEnvelope[T], error) {
	var envelope DebeziumEnvelope[T]
	err := json.Unmarshal(value, &envelope)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling message: %w", err)
	}

	return &envelope, nil
}

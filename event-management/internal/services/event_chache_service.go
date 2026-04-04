package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DevVictor19/event/internal/entities"
	"github.com/redis/go-redis/v9"
)

type EventCacheService interface {
	SetByUUID(ctx context.Context, event *entities.Event) error
	GetByUUID(ctx context.Context, uuid string) (*entities.Event, error)
}

type eventCacheService struct {
	rdb *redis.Client
}

func NewEventCacheService(rdb *redis.Client) EventCacheService {
	return &eventCacheService{
		rdb: rdb,
	}
}

func (s *eventCacheService) SetByUUID(ctx context.Context, event *entities.Event) error {
	oneHour := time.Hour
	jsonData, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return s.rdb.Set(ctx, s.getKey(event.UUID), jsonData, oneHour).Err()
}

func (s *eventCacheService) GetByUUID(ctx context.Context, uuid string) (*entities.Event, error) {
	data, err := s.rdb.Get(ctx, s.getKey(uuid)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var event entities.Event
	err = json.Unmarshal(data, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal event from cache: %w", err)
	}

	return &event, nil
}

func (s *eventCacheService) getKey(uuid string) string {
	return "event:" + uuid
}

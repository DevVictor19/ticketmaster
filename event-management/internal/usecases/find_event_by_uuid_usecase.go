package usecases

import (
	"context"

	"github.com/DevVictor19/event/internal/entities"
	"github.com/DevVictor19/event/internal/repositories"
	"github.com/DevVictor19/event/internal/services"
)

type FindEventByUuidUC struct {
	eventRepo     repositories.EventRepository
	eventCacheSvc services.EventCacheService
}

func NewFindEventByUuidUC(
	eventRepo repositories.EventRepository,
	eventCacheSvc services.EventCacheService) *FindEventByUuidUC {

	return &FindEventByUuidUC{
		eventRepo:     eventRepo,
		eventCacheSvc: eventCacheSvc,
	}
}

func (uc *FindEventByUuidUC) Execute(ctx context.Context, uuid string) (*entities.Event, error) {
	cached, err := uc.eventCacheSvc.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if cached != nil {
		return cached, nil
	}

	event, err := uc.eventRepo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	err = uc.eventCacheSvc.SetByUUID(ctx, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

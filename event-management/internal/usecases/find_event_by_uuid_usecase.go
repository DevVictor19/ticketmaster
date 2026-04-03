package usecases

import (
	"context"

	"github.com/DevVictor19/event/internal/entities"
	"github.com/DevVictor19/event/internal/repositories"
)

type FindEventByUuidUC struct {
	eventRepository repositories.EventRepository
}

func NewFindEventByUuidUC(eventRepository repositories.EventRepository) *FindEventByUuidUC {
	return &FindEventByUuidUC{
		eventRepository: eventRepository,
	}
}

func (uc *FindEventByUuidUC) Execute(ctx context.Context, uuid string) (*entities.Event, error) {
	return uc.eventRepository.FindByUUID(ctx, uuid)
}

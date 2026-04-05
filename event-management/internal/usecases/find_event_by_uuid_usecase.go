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
	ticketLockSvc services.TicketLockService
}

func NewFindEventByUuidUC(
	eventRepo repositories.EventRepository,
	eventCacheSvc services.EventCacheService,
	ticketLockSvc services.TicketLockService,
) *FindEventByUuidUC {

	return &FindEventByUuidUC{
		eventRepo:     eventRepo,
		eventCacheSvc: eventCacheSvc,
		ticketLockSvc: ticketLockSvc,
	}
}

func (uc *FindEventByUuidUC) Execute(ctx context.Context, uuid string) (*entities.Event, error) {
	var event *entities.Event

	cachedEv, err := uc.eventCacheSvc.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if cachedEv != nil {
		event = cachedEv
	} else {
		dbEvent, err := uc.eventRepo.FindByUUID(ctx, uuid)
		if err != nil {
			return nil, err
		}

		err = uc.eventCacheSvc.SetByUUID(ctx, dbEvent)
		if err != nil {
			return nil, err
		}

		event = dbEvent
	}

	reservations, err := uc.ticketLockSvc.CheckReservations(ctx, event.GetTicketUUIDs())
	if err != nil {
		return nil, err
	}

	event.UpdateReservationStatus(reservations)

	return event, nil
}

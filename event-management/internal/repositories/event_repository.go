package repositories

import (
	"context"

	"github.com/DevVictor19/event/internal/entities"
	"gorm.io/gorm"
)

type EventRepository interface {
	repository[entities.Event]
}

type eventRepository struct {
	baseRepository[entities.Event]
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{
		baseRepository: newBaseRepository[entities.Event](db),
	}
}

func (r *eventRepository) FindByUUID(ctx context.Context, uuid string) (*entities.Event, error) {
	entity, err := gorm.G[entities.Event](r.db).
		Where("uuid = ?", uuid).
		Preload("Venue", nil).
		Preload("Performers", nil).
		Preload("Tickets", nil).
		First(ctx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &entity, nil
}

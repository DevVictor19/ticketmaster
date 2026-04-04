package repositories

import (
	"context"

	"github.com/DevVictor19/booking/internal/entities"
	"gorm.io/gorm"
)

type TicketRepository interface {
	repository[entities.Ticket]
	CheckTicketAvailability(ctx context.Context, uuid string) (bool, error)
}

type ticketRepository struct {
	baseRepository[entities.Ticket]
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{
		baseRepository: newBaseRepository[entities.Ticket](db),
	}
}

func (r *ticketRepository) CheckTicketAvailability(ctx context.Context, uuid string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Ticket{}).
		Where("uuid = ? AND status = ?", uuid, entities.TicketStatusAvailable).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

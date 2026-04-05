package services

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type TicketLockService interface {
	// GetReservations returns a slice of booleans indicating whether each ticket UUID is reserved (true) or not (false).
	GetReservations(ctx context.Context, ticketUUIDs []string) ([]bool, error)
}

type ticketLockService struct {
	rdb *redis.Client
}

func NewTicketLockService(rdb *redis.Client) TicketLockService {
	return &ticketLockService{
		rdb: rdb,
	}
}

func (s *ticketLockService) GetReservations(ctx context.Context, ticketUUIDs []string) ([]bool, error) {
	keys := make([]string, len(ticketUUIDs))
	for i, uuid := range ticketUUIDs {
		keys[i] = s.getLockKey(uuid)
	}

	results, err := s.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	reservations := make([]bool, len(results))
	fmt.Println(reservations)
	for i, res := range results {
		if res == nil {
			reservations[i] = false
		} else {
			reservations[i] = true
		}
	}

	return reservations, nil
}

func (s *ticketLockService) getLockKey(ticketUUID string) string {
	return "ticket_lock:" + ticketUUID
}

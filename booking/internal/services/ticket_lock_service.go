package services

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrTicketAlreadyReserved = errors.New("ticket is already reserved")

type TicketLockService interface {
	// ReserveTicket attempts to reserve a ticket by its UUID. If the ticket is successfully reserved, it will be locked for the specified TTL duration.
	ReserveTicket(ctx context.Context, ticketUUID string, userUUID string) error
}

type ticketLockService struct {
	rdb *redis.Client
}

func NewTicketLockService(rdb *redis.Client) TicketLockService {
	return &ticketLockService{
		rdb: rdb,
	}
}

func (s *ticketLockService) ReserveTicket(
	ctx context.Context,
	ticketUUID string,
	userUUID string) error {

	lockTTL := time.Minute * 15
	pipe := s.rdb.Pipeline()

	pipe.SetArgs(ctx, s.getLockKey(ticketUUID), "true", redis.SetArgs{
		Mode: "NX",
		TTL:  lockTTL,
	})
	pipe.SetArgs(ctx, s.getUserLockKey(ticketUUID, userUUID), "true", redis.SetArgs{
		Mode: "NX",
		TTL:  lockTTL,
	})

	_, err := pipe.Exec(ctx)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrTicketAlreadyReserved
		}
		return err
	}

	return nil
}

func (s *ticketLockService) getLockKey(ticketUUID string) string {
	return "ticket_lock:" + ticketUUID
}

func (s *ticketLockService) getUserLockKey(ticketUUID string, userUUID string) string {
	return "ticket_lock_user" + ticketUUID + ":" + userUUID
}

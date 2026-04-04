package usecases

import (
	"context"
	"errors"

	"github.com/DevVictor19/booking/internal/repositories"
	"github.com/DevVictor19/booking/internal/services"
)

type ReserveTicketUseCase struct {
	ticketRepo    repositories.TicketRepository
	ticketLockSvc services.TicketLockService
}

func NewReserveTicketUseCase(
	ticketRepo repositories.TicketRepository,
	ticketLockSvc services.TicketLockService) *ReserveTicketUseCase {

	return &ReserveTicketUseCase{
		ticketRepo:    ticketRepo,
		ticketLockSvc: ticketLockSvc,
	}
}

var ErrTicketNotAvailable = errors.New("ticket not available for reservation")

type ReserveTicketInput struct {
	TicketUUID string
	UserUUID   string
}

func (uc *ReserveTicketUseCase) Execute(ctx context.Context, input ReserveTicketInput) error {
	isAvailable, err := uc.ticketRepo.CheckAvailability(ctx, input.TicketUUID)
	if err != nil {
		return err
	}
	if !isAvailable {
		return ErrTicketNotAvailable
	}

	err = uc.ticketLockSvc.ReserveTicket(ctx, input.TicketUUID, input.UserUUID)
	if err != nil {
		if errors.Is(err, services.ErrTicketAlreadyReserved) {
			return ErrTicketNotAvailable
		}
		return err
	}

	return nil
}

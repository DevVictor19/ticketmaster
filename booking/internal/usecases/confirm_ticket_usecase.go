package usecases

import (
	"context"
	"errors"

	"github.com/DevVictor19/booking/internal/entities"
	"github.com/DevVictor19/booking/internal/repositories"
	"github.com/DevVictor19/booking/internal/services"
)

type ConfirmTicketUseCase struct {
	ticketRepo    repositories.TicketRepository
	ticketLockSvc services.TicketLockService
}

func NewConfirmTicketUseCase(
	ticketRepo repositories.TicketRepository,
	ticketLockSvc services.TicketLockService) *ConfirmTicketUseCase {

	return &ConfirmTicketUseCase{
		ticketRepo:    ticketRepo,
		ticketLockSvc: ticketLockSvc,
	}
}

var ErrTicketNotReserved = errors.New("ticket is not reserved for this user")

type ConfirmTicketInput struct {
	TicketUUID        string
	UserUUID          string
	PaymentMethodUUID string
}

func (uc *ConfirmTicketUseCase) Execute(ctx context.Context, input ConfirmTicketInput) error {
	isReserved, err := uc.ticketLockSvc.CheckReservation(ctx, input.TicketUUID, input.UserUUID)
	if err != nil {
		return err
	}
	if !isReserved {
		return ErrTicketNotReserved
	}

	// TODO: Implement payment integration here
	// TOdo: publish event on Kafka to notify other services about the confirmation

	return uc.ticketRepo.UpdateStatus(ctx, input.TicketUUID, entities.TicketStatusSold)
}

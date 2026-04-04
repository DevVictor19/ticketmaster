package handlers

import (
	"github.com/DevVictor19/booking/internal/handlers/dtos"
	"github.com/DevVictor19/booking/internal/usecases"
	"github.com/DevVictor19/booking/internal/utils"
	"github.com/labstack/echo/v5"
)

type TicketsHandler struct {
	reserveTicketUC *usecases.ReserveTicketUseCase
	confirmTicketUC *usecases.ConfirmTicketUseCase
}

func NewTicketsHandler(
	reserveTicketUC *usecases.ReserveTicketUseCase,
	confirmTicketUC *usecases.ConfirmTicketUseCase) *TicketsHandler {

	return &TicketsHandler{
		reserveTicketUC: reserveTicketUC,
		confirmTicketUC: confirmTicketUC,
	}
}

func (h *TicketsHandler) ReserveTicket(c *echo.Context) error {
	var req dtos.ReserveTicketRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, utils.ErrorResponse(err, "Invalid request payload"))
	}

	input := usecases.ReserveTicketInput{
		TicketUUID: req.TicketUUID,
		UserUUID:   req.UserUUID,
	}

	err := h.reserveTicketUC.Execute(c.Request().Context(), input)
	if err != nil {
		if err == usecases.ErrTicketNotAvailable {
			return c.JSON(409, utils.ErrorResponse(err, "Ticket is not available for reservation"))
		}
		return c.JSON(500, utils.ErrorResponse(err, "Internal server error"))
	}

	return c.NoContent(200)
}

func (h *TicketsHandler) ConfirmTicket(c *echo.Context) error {
	var req dtos.ConfirmTicketRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, utils.ErrorResponse(err, "Invalid request payload"))
	}

	input := usecases.ConfirmTicketInput{
		TicketUUID:        req.TicketUUID,
		UserUUID:          req.UserUUID,
		PaymentMethodUUID: req.PaymentMethodUUID,
	}

	err := h.confirmTicketUC.Execute(c.Request().Context(), input)
	if err != nil {
		if err == usecases.ErrTicketNotReserved {
			return c.JSON(409, utils.ErrorResponse(err, "Ticket is not reserved for this user"))
		}
		return c.JSON(500, utils.ErrorResponse(err, "Internal server error"))
	}

	return c.NoContent(200)
}

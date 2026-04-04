package handlers

import (
	"github.com/DevVictor19/booking/internal/handlers/dtos"
	"github.com/DevVictor19/booking/internal/usecases"
	"github.com/DevVictor19/booking/internal/utils"
	"github.com/labstack/echo/v5"
)

type TicketsHandler struct {
	reserveTicketUC *usecases.ReserveTicketUseCase
}

func NewTicketsHandler(reserveTicketUC *usecases.ReserveTicketUseCase) *TicketsHandler {
	return &TicketsHandler{
		reserveTicketUC: reserveTicketUC,
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

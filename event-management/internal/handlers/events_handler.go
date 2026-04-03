package handlers

import (
	"errors"

	"github.com/DevVictor19/event/internal/handlers/dtos"
	"github.com/DevVictor19/event/internal/repositories"
	"github.com/DevVictor19/event/internal/usecases"
	"github.com/DevVictor19/event/internal/utils"
	"github.com/labstack/echo/v5"
)

type EventsHandler struct {
	findEventByUuidUC *usecases.FindEventByUuidUC
}

func NewEventsHandler(findEventByUuidUC *usecases.FindEventByUuidUC) *EventsHandler {
	return &EventsHandler{
		findEventByUuidUC: findEventByUuidUC,
	}
}

func (h *EventsHandler) FindByUUID(c *echo.Context) error {
	uuid := c.Param("uuid")
	event, err := h.findEventByUuidUC.Execute(c.Request().Context(), uuid)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return c.JSON(404, utils.ErrorResponse(err, "Event not found"))
		}
		return c.JSON(500, utils.ErrorResponse(err, "Internal server error"))
	}
	return c.JSON(200, dtos.ToEventDTO(event))
}

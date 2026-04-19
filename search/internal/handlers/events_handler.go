package handlers

import (
	"strconv"

	"github.com/DevVictor19/search/internal/handlers/dtos"
	"github.com/DevVictor19/search/internal/usecases"
	"github.com/DevVictor19/search/internal/utils"
	"github.com/labstack/echo/v5"
)

type EventsHandler struct {
	searchEventsUC *usecases.SearchEventsUC
}

func NewEventsHandler(searchEventsUC *usecases.SearchEventsUC) *EventsHandler {
	return &EventsHandler{
		searchEventsUC: searchEventsUC,
	}
}

func (h *EventsHandler) SearchEvents(c *echo.Context) error {
	q := c.QueryParam("q")

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	size, err := strconv.Atoi(c.QueryParam("size"))
	if err != nil || size <= 0 || size > 100 {
		size = 10
	}

	result, err := h.searchEventsUC.Execute(c.Request().Context(), usecases.SearchEventsInput{
		Query: q,
		Page:  page,
		Size:  size,
	})
	if err != nil {
		return c.JSON(500, utils.ErrorResponse(err, "Internal server error"))
	}

	return c.JSON(200, dtos.ToSearchEventsResponseDTO(result))
}

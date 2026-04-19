package dtos

import (
	"time"

	"github.com/DevVictor19/search/internal/services"
)

type EventDocDTO struct {
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
}

type SearchEventsResponseDTO struct {
	Events []EventDocDTO `json:"events"`
	Total  int           `json:"total"`
	Page   int           `json:"page"`
	Size   int           `json:"size"`
}

func ToSearchEventsResponseDTO(res *services.SearchEventsResponse) *SearchEventsResponseDTO {
	events := make([]EventDocDTO, len(res.Events))
	for i, e := range res.Events {
		events[i] = EventDocDTO{
			UUID:        e.UUID,
			Name:        e.Name,
			Description: e.Description,
			Date:        e.Date,
			Location:    e.Location,
		}
	}

	return &SearchEventsResponseDTO{
		Events: events,
		Total:  res.Total,
		Page:   res.Page,
		Size:   res.Size,
	}
}

package dtos

import (
	"time"

	"github.com/DevVictor19/event/internal/entities"
)

type EventDTO struct {
	UUID        string          `json:"uuid"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Date        time.Time       `json:"date"`
	Venue       *VenueDTO       `json:"venue,omitempty"`
	Performers  []*PerformerDTO `json:"performers,omitempty"`
	Tickets     []*TicketDTO    `json:"tickets,omitempty"`
}

func ToEventDTO(event *entities.Event) *EventDTO {
	if event == nil {
		return nil
	}

	dto := &EventDTO{
		UUID:        event.UUID,
		Name:        event.Name,
		Description: event.Description,
		Date:        event.Date,
	}

	if event.Venue != nil {
		dto.Venue = ToVenueDTO(event.Venue)
	}

	if len(event.Performers) > 0 {
		dto.Performers = make([]*PerformerDTO, len(event.Performers))
		for i, performer := range event.Performers {
			dto.Performers[i] = ToPerformerDTO(performer)
		}
	}

	if len(event.Tickets) > 0 {
		dto.Tickets = make([]*TicketDTO, len(event.Tickets))
		for i, ticket := range event.Tickets {
			dto.Tickets[i] = ToTicketDTO(ticket)
		}
	}

	return dto
}

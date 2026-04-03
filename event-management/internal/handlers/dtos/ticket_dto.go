package dtos

import "github.com/DevVictor19/event/internal/entities"

type TicketDTO struct {
	UUID   string    `json:"uuid"`
	Event  *EventDTO `json:"event,omitempty"`
	Price  uint      `json:"price"`
	Seat   string    `json:"seat"`
	Status string    `json:"status"`
}

func ToTicketDTO(ticket *entities.Ticket) *TicketDTO {
	if ticket == nil {
		return nil
	}

	var eventDTO *EventDTO
	if ticket.Event != nil {
		eventDTO = ToEventDTO(ticket.Event)
	}

	return &TicketDTO{
		UUID:   ticket.UUID,
		Event:  eventDTO,
		Price:  ticket.Price,
		Seat:   ticket.Seat,
		Status: string(ticket.Status),
	}
}

package entities

import "time"

type Event struct {
	BaseEntity
	Venue       *Venue       `json:"venue,omitempty"`
	VenueID     uint         `json:"venue_id"`
	Performers  []*Performer `json:"performers,omitempty" gorm:"many2many:event_performers"`
	Tickets     []*Ticket    `json:"tickets,omitempty"`
	Date        time.Time    `json:"date" gorm:"index"`
	Name        string       `json:"name" gorm:"index"`
	Description *string      `json:"description" gorm:"type:text"`
}

func (e *Event) GetTicketUUIDs() []string {
	uuids := make([]string, len(e.Tickets))
	for i, ticket := range e.Tickets {
		uuids[i] = ticket.UUID
	}
	return uuids
}

func (e *Event) UpdateReservedTickets(reservations []bool) {
	for i, reserved := range reservations {
		if reserved {
			e.Tickets[i].Status = TicketStatusReserved
		}
	}
}

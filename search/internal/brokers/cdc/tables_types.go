package cdc

import "time"

type TicketStatus string

const (
	TicketStatusAvailable TicketStatus = "available"
	TicketStatusBooked    TicketStatus = "booked"
)

type Ticket struct {
	ID        uint         `json:"id"`
	UUID      string       `json:"uuid"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	EventID   uint         `json:"event_id"`
	Price     uint         `json:"price"`
	Seat      string       `json:"seat"`
	Status    TicketStatus `json:"status"`
}

type Event struct {
	ID          uint      `json:"id"`
	UUID        string    `json:"uuid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Date        time.Time `json:"date"`
	VenueID     uint      `json:"venue_id"`
}

type Venue struct {
	ID        uint            `json:"id"`
	UUID      string          `json:"uuid"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Location  string          `json:"location"`
	SeatMap   map[string]bool `json:"seat_map"`
}

type Performer struct {
	ID          uint      `json:"id"`
	UUID        string    `json:"uuid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Age         uint      `json:"age"`
	Description *string   `json:"description"`
}

type EventPerformer struct {
	EventID     uint `json:"event_id"`
	PerformerID uint `json:"performer_id"`
}

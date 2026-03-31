package entities

type Ticket struct {
	baseEntity
	EventID uint         `json:"event_id" gorm:"index"`
	Event   Event        `json:"event" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Price   uint         `json:"price"`
	Seat    string       `json:"seat" gorm:"index"`
	Status  TicketStatus `json:"status" gorm:"index"`
}

type TicketStatus string

const (
	TicketStatusAvailable TicketStatus = "available"
	TicketStatusSold      TicketStatus = "booked"
)

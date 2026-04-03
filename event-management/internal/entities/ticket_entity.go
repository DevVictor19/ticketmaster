package entities

type Ticket struct {
	BaseEntity
	EventID uint         `json:"event_id" gorm:"index"`
	Event   *Event       `json:"event,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Price   uint         `json:"price"`
	Seat    string       `json:"seat" gorm:"index"`
	Status  TicketStatus `json:"status,omitempty" gorm:"index"`
}

type TicketStatus string

const (
	TicketStatusAvailable TicketStatus = "available"
	TicketStatusSold      TicketStatus = "booked"
)

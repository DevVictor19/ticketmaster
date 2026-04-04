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

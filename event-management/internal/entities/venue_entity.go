package entities

type Venue struct {
	BaseEntity
	Events   []*Event        `json:"events" gorm:"foreignKey:VenueID"`
	Location string          `json:"location" gorm:"index"`
	SeatMap  map[string]bool `json:"seat_map" gorm:"type:json"`
}

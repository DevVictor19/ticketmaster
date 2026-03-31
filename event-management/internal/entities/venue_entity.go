package entities

type Venue struct {
	baseEntity
	Events   []Event         `json:"events"`
	Location string          `json:"location" gorm:"index"`
	SeatMap  map[string]bool `json:"seat_map" gorm:"type:json"`
}

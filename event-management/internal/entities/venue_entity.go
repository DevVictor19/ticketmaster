package entities

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Venue struct {
	BaseEntity
	Events   []*Event `json:"events,omitempty" gorm:"foreignKey:VenueID"`
	Location string   `json:"location" gorm:"index"`
	SeatMap  SeatMap  `json:"seat_map" gorm:"type:json"`
}

type SeatMap map[string]bool

func (s SeatMap) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SeatMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan SeatMap: %v", value)
	}
	return json.Unmarshal(bytes, s)
}

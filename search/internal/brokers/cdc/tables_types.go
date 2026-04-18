package cdc

import (
	"encoding/json"
	"time"
)

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
	ID        uint      `json:"id"`
	UUID      string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Location  string    `json:"location"`
	SeatMap   SeatMap   `json:"seat_map"`
}

type SeatMap map[string]bool

func (s *SeatMap) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		return json.Unmarshal([]byte(str), (*map[string]bool)(s))
	}
	return json.Unmarshal(data, (*map[string]bool)(s))
}

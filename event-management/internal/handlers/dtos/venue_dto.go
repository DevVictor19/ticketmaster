package dtos

import "github.com/DevVictor19/event/internal/entities"

type VenueDTO struct {
	UUID     string          `json:"uuid"`
	Location string          `json:"location"`
	SeatMap  map[string]bool `json:"seat_map"`
	Events   []*EventDTO     `json:"events,omitempty"`
}

func ToVenueDTO(venue *entities.Venue) *VenueDTO {
	if venue == nil {
		return nil
	}

	dto := &VenueDTO{
		UUID:     venue.UUID,
		Location: venue.Location,
		SeatMap:  venue.SeatMap,
	}

	if len(venue.Events) > 0 {
		dto.Events = make([]*EventDTO, len(venue.Events))
		for i, event := range venue.Events {
			dto.Events[i] = ToEventDTO(event)
		}
	}

	return dto
}

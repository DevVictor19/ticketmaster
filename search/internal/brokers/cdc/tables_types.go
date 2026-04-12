package cdc

type Ticket struct {
	ID        int64  `json:"id"`
	UUID      string `json:"uuid"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	EventID   int64  `json:"event_id"`
	Price     int64  `json:"price"`
	Seat      string `json:"seat"`
	Status    string `json:"status"`
}

type Event struct {
	ID          int64  `json:"id"`
	UUID        string `json:"uuid"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Name        string `json:"name"`
	Description string `json:"description"`
	EventDate   string `json:"event_date"`
	VenueID     int64  `json:"venue_id"`
}

type Venue struct {
	ID        int64  `json:"id"`
	UUID      string `json:"uuid"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	Capacity  int64  `json:"capacity"`
}

type Performer struct {
	ID        int64  `json:"id"`
	UUID      string `json:"uuid"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	Genre     string `json:"genre"`
}

type EventPerformer struct {
	EventID     int64 `json:"event_id"`
	PerformerID int64 `json:"performer_id"`
}

package dtos

type ReserveTicketRequest struct {
	TicketUUID string `json:"ticket_uuid" validate:"required,uuid4"`
	UserUUID   string `json:"user_uuid" validate:"required,uuid4"`
}

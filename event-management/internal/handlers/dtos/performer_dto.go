package dtos

import "github.com/DevVictor19/event/internal/entities"

type PerformerDTO struct {
	UUID        string      `json:"uuid"`
	Name        string      `json:"name"`
	Age         uint        `json:"age"`
	Description *string     `json:"description,omitempty"`
	Events      []*EventDTO `json:"events,omitempty"`
}

func ToPerformerDTO(performer *entities.Performer) *PerformerDTO {
	if performer == nil {
		return nil
	}

	dto := &PerformerDTO{
		UUID:        performer.UUID,
		Name:        performer.Name,
		Age:         performer.Age,
		Description: performer.Description,
	}

	if len(performer.Events) > 0 {
		dto.Events = make([]*EventDTO, len(performer.Events))
		for i, event := range performer.Events {
			dto.Events[i] = ToEventDTO(event)
		}
	}

	return dto
}

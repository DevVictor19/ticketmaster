package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/DevVictor19/search/internal/entities"
	"github.com/elastic/go-elasticsearch/v9"
)

const eventIndex = "events"

type SearchEngineService interface {
	UpsertEventIdx(event *entities.Event) error
	DeleteEventIdx(eventID uint) error
	UpdateVenueIdx(venue *entities.Venue) error
}

type elasticsearchService struct {
	client *elasticsearch.Client
}

func NewSearchEngineService(client *elasticsearch.Client) SearchEngineService {
	return &elasticsearchService{client: client}
}

type EventDoc struct {
	ID            uint      `json:"event_id"`
	UUID          string    `json:"event_uuid"`
	Name          string    `json:"event_name"`
	Description   *string   `json:"event_description,omitempty"`
	Date          time.Time `json:"event_date"`
	VenueID       uint      `json:"venue_id"`
	VenueUUID     string    `json:"venue_uuid"`
	VenueLocation string    `json:"venue_location"`
}

func (s *elasticsearchService) UpsertEventIdx(event *entities.Event) error {
	doc := EventDoc{
		ID:          event.ID,
		UUID:        event.UUID,
		Name:        event.Name,
		Description: event.Description,
		Date:        event.Date,
		VenueID:     event.VenueID,
	}

	if event.Venue != nil {
		doc.VenueUUID = event.Venue.UUID
		doc.VenueLocation = event.Venue.Location
	}

	// TODO: populate VenueUUID and VenueLocation when Venue is nil (fetch from store)

	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal event document: %w", err)
	}

	docID := strconv.FormatUint(uint64(event.ID), 10)

	res, err := s.client.Index(
		eventIndex,
		bytes.NewReader(body),
		s.client.Index.WithDocumentID(docID),
		s.client.Index.WithContext(context.Background()),
	)
	if err != nil {
		return fmt.Errorf("failed to index event: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch index error: %s", res.String())
	}

	slog.Info("event indexed", "id", event.ID)
	return nil
}

func (s *elasticsearchService) DeleteEventIdx(eventID uint) error {
	docID := strconv.FormatUint(uint64(eventID), 10)

	res, err := s.client.Delete(
		eventIndex,
		docID,
		s.client.Delete.WithContext(context.Background()),
	)
	if err != nil {
		return fmt.Errorf("failed to delete event from index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch delete error: %s", res.String())
	}

	slog.Info("event deleted from index", "id", eventID)
	return nil
}

func (s *elasticsearchService) UpdateVenueIdx(venue *entities.Venue) error {
	query := map[string]interface{}{
		"script": map[string]interface{}{
			"source": `ctx._source.venue_uuid = params.venue_uuid;
						ctx._source.venue_location = params.venue_location;`,
			"params": map[string]interface{}{
				"venue_uuid":     venue.UUID,
				"venue_location": venue.Location,
			},
		},
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"venue_id": venue.ID,
			},
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		return fmt.Errorf("failed to marshal update query: %w", err)
	}

	res, err := s.client.UpdateByQuery(
		[]string{eventIndex},
		s.client.UpdateByQuery.WithBody(bytes.NewReader(body)),
		s.client.UpdateByQuery.WithContext(context.Background()),
	)
	if err != nil {
		return fmt.Errorf("failed to update venue in index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch update_by_query error: %s", res.String())
	}

	slog.Info("venue updated in index", "id", venue.ID)
	return nil
}

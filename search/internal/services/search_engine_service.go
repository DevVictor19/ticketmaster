package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/DevVictor19/search/internal/repositories"
	"github.com/elastic/go-elasticsearch/v9"
)

const eventIndex = "events"

type SearchEngineService interface {
	UpsertEventIdx(doc *EventDoc) error
	DeleteEventIdx(eventID uint) error
	UpdateEventLocationIdx(venueID uint, location string) error
}

type elasticsearchService struct {
	client    *elasticsearch.Client
	venueRepo repositories.VenueRepository
}

func NewSearchEngineService(
	client *elasticsearch.Client,
	venueRepo repositories.VenueRepository,
) SearchEngineService {

	return &elasticsearchService{
		client:    client,
		venueRepo: venueRepo,
	}
}

type EventDoc struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Date        time.Time `json:"date"`
	VenueID     uint      `json:"venue_id"`
	Location    string    `json:"location"`
}

func (s *elasticsearchService) UpsertEventIdx(doc *EventDoc) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	venue, err := s.venueRepo.FindByID(ctx, doc.VenueID)
	if err != nil {
		return fmt.Errorf("failed to fetch venue for event: %w", err)
	}
	doc.Location = venue.Location

	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal event document: %w", err)
	}

	docID := strconv.FormatUint(uint64(doc.ID), 10)

	res, err := s.client.Index(
		eventIndex,
		bytes.NewReader(body),
		s.client.Index.WithDocumentID(docID),
		s.client.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to index event: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch index error: %s", res.String())
	}

	slog.Info("event indexed", "id", doc.ID)
	return nil
}

func (s *elasticsearchService) DeleteEventIdx(eventID uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	docID := strconv.FormatUint(uint64(eventID), 10)

	res, err := s.client.Delete(
		eventIndex,
		docID,
		s.client.Delete.WithContext(ctx),
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

func (s *elasticsearchService) UpdateEventLocationIdx(venueID uint, location string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := map[string]interface{}{
		"script": map[string]interface{}{
			"source": `ctx._source.location = params.location;`,
			"params": map[string]interface{}{
				"location": location,
			},
		},
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"venue_id": venueID,
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
		s.client.UpdateByQuery.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to update events location: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch update_by_query error: %s", res.String())
	}

	slog.Info("updated events location", "venue_id", venueID)
	return nil
}

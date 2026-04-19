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
	UpsertEventIdx(ctx context.Context, doc *EventDoc) error
	DeleteEventIdx(ctx context.Context, eventID uint) error
	UpdateEventLocationIdx(ctx context.Context, venueID uint, location string) error
	SearchEvents(ctx context.Context, req SearchEventsRequest) (*SearchEventsResponse, error)
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

type SearchEventsRequest struct {
	Query string
	Page  int
	Size  int
}

type SearchEventsResponse struct {
	Events []EventDoc
	Total  int
	Page   int
	Size   int
}

func (s *elasticsearchService) UpsertEventIdx(ctx context.Context, doc *EventDoc) error {
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

func (s *elasticsearchService) DeleteEventIdx(ctx context.Context, eventID uint) error {
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

func (s *elasticsearchService) UpdateEventLocationIdx(ctx context.Context, venueID uint, location string) error {
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

func (s *elasticsearchService) SearchEvents(ctx context.Context, req SearchEventsRequest) (*SearchEventsResponse, error) {
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	from := (req.Page - 1) * req.Size

	query := map[string]interface{}{
		"from": from,
		"size": req.Size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"query":     req.Query,
							"fields":    []string{"name^4", "description", "location^2"},
							"fuzziness": "AUTO",
						},
					},
				},
				"filter": []interface{}{
					map[string]interface{}{
						"range": map[string]interface{}{
							"date": map[string]interface{}{
								"gte": "now",
							},
						},
					},
				},
			},
		},
		"sort": []interface{}{
			map[string]interface{}{
				"_score": "desc",
			},
			map[string]interface{}{
				"date": "asc",
			},
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(eventIndex),
		s.client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var r struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source EventDoc `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	events := make([]EventDoc, 0, len(r.Hits.Hits))
	for _, hit := range r.Hits.Hits {
		events = append(events, hit.Source)
	}

	return &SearchEventsResponse{
		Events: events,
		Total:  r.Hits.Total.Value,
		Page:   req.Page,
		Size:   req.Size,
	}, nil
}

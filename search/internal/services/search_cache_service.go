package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type SearchCacheService interface {
	Set(ctx context.Context, req SearchEventsRequest, res *SearchEventsResponse) error
	Get(ctx context.Context, req SearchEventsRequest) (*SearchEventsResponse, error)
}

type searchCacheService struct {
	rdb *redis.Client
}

func NewSearchCacheService(rdb *redis.Client) SearchCacheService {
	return &searchCacheService{
		rdb: rdb,
	}
}

func (s *searchCacheService) Set(ctx context.Context, req SearchEventsRequest, res *SearchEventsResponse) error {
	jsonData, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal search response: %w", err)
	}
	return s.rdb.Set(ctx, s.getKey(req), jsonData, 5*time.Minute).Err()
}

func (s *searchCacheService) Get(ctx context.Context, req SearchEventsRequest) (*SearchEventsResponse, error) {
	data, err := s.rdb.Get(ctx, s.getKey(req)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var res SearchEventsResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response from cache: %w", err)
	}

	return &res, nil
}

func (s *searchCacheService) getKey(req SearchEventsRequest) string {
	q := strings.TrimSpace(strings.ToLower(req.Query))
	return fmt.Sprintf("search:%s:%d:%d", q, req.Page, req.Size)
}

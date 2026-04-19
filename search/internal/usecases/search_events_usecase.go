package usecases

import (
	"context"

	"github.com/DevVictor19/search/internal/services"
)

type SearchEventsInput struct {
	Query string
	Page  int
	Size  int
}

type SearchEventsUC struct {
	searchEngineSvc services.SearchEngineService
	searchCacheSvc  services.SearchCacheService
}

func NewSearchEventsUC(
	searchEngineSvc services.SearchEngineService,
	searchCacheSvc services.SearchCacheService,
) *SearchEventsUC {

	return &SearchEventsUC{
		searchEngineSvc: searchEngineSvc,
		searchCacheSvc:  searchCacheSvc,
	}
}

func (uc *SearchEventsUC) Execute(ctx context.Context, input SearchEventsInput) (*services.SearchEventsResponse, error) {
	req := services.SearchEventsRequest{
		Query: input.Query,
		Page:  input.Page,
		Size:  input.Size,
	}

	cached, err := uc.searchCacheSvc.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	if cached != nil {
		return cached, nil
	}

	result, err := uc.searchEngineSvc.SearchEvents(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := uc.searchCacheSvc.Set(ctx, req, result); err != nil {
		return nil, err
	}

	return result, nil
}

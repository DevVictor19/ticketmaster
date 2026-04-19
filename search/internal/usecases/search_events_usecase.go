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
}

func NewSearchEventsUC(searchEngineSvc services.SearchEngineService) *SearchEventsUC {
	return &SearchEventsUC{
		searchEngineSvc: searchEngineSvc,
	}
}

func (uc *SearchEventsUC) Execute(ctx context.Context, input SearchEventsInput) (*services.SearchEventsResponse, error) {
	return uc.searchEngineSvc.SearchEvents(ctx, services.SearchEventsRequest{
		Query: input.Query,
		Page:  input.Page,
		Size:  input.Size,
	})
}

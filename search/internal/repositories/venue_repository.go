package repositories

import (
	"github.com/DevVictor19/search/internal/entities"
	"gorm.io/gorm"
)

type VenueRepository interface {
	repository[entities.Venue]
}

type venueRepository struct {
	baseRepository[entities.Venue]
}

func NewVenueRepository(db *gorm.DB) VenueRepository {
	return &venueRepository{
		baseRepository: newBaseRepository[entities.Venue](db),
	}
}

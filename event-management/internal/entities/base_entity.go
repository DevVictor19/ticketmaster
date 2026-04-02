package entities

import (
	"time"

	"github.com/google/uuid"
)

type baseEntity struct {
	ID        uint      `json:"id"`
	UUID      string    `json:"uuid" gorm:"uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *baseEntity) BeforeCreate() {
	b.UUID = uuid.New().String()
}

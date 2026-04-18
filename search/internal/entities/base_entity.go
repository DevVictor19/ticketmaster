package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UUID      string    `json:"uuid" gorm:"uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *BaseEntity) BeforeCreate(db *gorm.DB) error {
	b.UUID = uuid.New().String()
	return nil
}

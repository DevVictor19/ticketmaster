package entities

import "time"

type baseEntity struct {
	ID        uint      `json:"id"`
	UUID      string    `json:"uuid" gorm:"uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

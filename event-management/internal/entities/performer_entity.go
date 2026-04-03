package entities

type Performer struct {
	BaseEntity
	Events      []*Event `json:"events" gorm:"many2many:event_performers"`
	Age         uint     `json:"age"`
	Name        string   `json:"name" gorm:"index"`
	Description *string  `json:"description" gorm:"type:text"`
}

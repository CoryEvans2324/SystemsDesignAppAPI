package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Hut struct {
	gorm.Model

	ObjectID   uint `gorm:"primaryKey"`
	GlobalID   uuid.UUID
	Place      string
	Region     string
	Status     string
	Bookable   bool
	Facilities string
	HasAlerts  string
	Thumbnail  string `gorm:"type:url"`
	StaticLink string `gorm:"type:url"`
	Location   Point
}

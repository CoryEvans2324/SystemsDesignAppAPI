package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Campsite struct {
	gorm.Model

	ObjectID       uint `gorm:"primaryKey"`
	GlobalID       uuid.UUID
	Name           string
	Region         string
	LocationString string
	Category       string
	Status         string
	Free           bool
	Facilities     string
	Activities     string
	DogsAllowed    bool
	Landscape      string
	HasAlerts      string
	Access         string
	Thumbnail      string `gorm:"type:url"`
	StaticLink     string `gorm:"type:url"`
	Location       Point
}

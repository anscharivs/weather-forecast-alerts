package models

import (
	"time"

	"gorm.io/gorm"
)

type Weather struct {
	gorm.Model
	CityID         uint
	Temperature    float64
	MinTemperature float64
	MaxTemperature float64
	Pressure       int
	Humidity       int
	Visibility     int
	Description    string
	FetchedAt      time.Time
}

package models

import "gorm.io/gorm"

type Weather struct {
	gorm.Model
	CityID      uint
	Temperature float64
}

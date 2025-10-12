package models

import "gorm.io/gorm"

type Alert struct {
	gorm.Model
	CityID  uint
	Type    string
	City    City
	Message string
}

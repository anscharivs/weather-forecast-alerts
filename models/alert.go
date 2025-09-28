package models

import "gorm.io/gorm"

type Alert struct {
	gorm.Model
	CityID  uint
	Message string
}

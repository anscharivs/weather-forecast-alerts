package models

import "gorm.io/gorm"

type City struct {
	gorm.Model
	Name string `gorm:"not null;size:256"` // https://gorm.io/docs/models.html#Fields-Tags
}

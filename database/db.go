package database

import (
	config "github.com/anscharivs/weather-forecast-alerts/internal"
	"github.com/anscharivs/weather-forecast-alerts/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.City{}, &models.Weather{}, &models.Alert{})

	return db, nil
}

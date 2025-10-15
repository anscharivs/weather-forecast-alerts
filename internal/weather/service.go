package weather

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	config "github.com/anscharivs/weather-forecast-alerts/internal"
	"github.com/anscharivs/weather-forecast-alerts/models"
	"gorm.io/gorm"
)

var BaseURL = "https://api.openweathermap.org/data/2.5/weather"

type apiResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
		Humidity int     `json:"humidity"`
		Pressure int     `json:"pressure"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Dt         int64 `json:"dt"`
	Visibility int   `json:"visibility"`
}

func FetchAndStoreWeatherData(db *gorm.DB, cfg config.Config) {

	var cities []models.City

	db.Find(&cities)

	for _, city := range cities {

		escapedCity := url.QueryEscape(city.Name)

		url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric&lang=es", BaseURL, escapedCity, cfg.APIKey)

		res, err := http.Get(url)

		if err != nil {
			fmt.Println("Error for ", city.Name, " weather ", err)
			continue
		}

		defer res.Body.Close()

		if res.StatusCode == 404 {
			continue
		}

		var data apiResponse

		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			fmt.Println("Response decoding error", err)
			continue
		}

		weather := models.Weather{
			CityID:         city.ID,
			Temperature:    data.Main.Temp,
			MinTemperature: data.Main.TempMin,
			MaxTemperature: data.Main.TempMax,
			Pressure:       data.Main.Pressure,
			Humidity:       data.Main.Humidity,
			Visibility:     data.Visibility,
			Description:    data.Weather[0].Description,
			Icon:           data.Weather[0].Icon,
			FetchedAt:      time.Unix(data.Dt, 0).UTC(),
		}

		db.Create(&weather)

		CheckForAlert(db, &weather)

		fmt.Printf("Saved %s: %.1f°C, %s\n", city.Name, weather.Temperature, weather.Description)
	}

}

func StartWeatherPolling(db *gorm.DB, cfg config.Config, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			FetchAndStoreWeatherData(db, cfg)
		}
	}()
}

func ExistsInDB(db *gorm.DB, name string) bool {
	var city models.City
	result := db.Where("name = ?", name).First(&city)
	return result.RowsAffected > 0
}

func CheckForAlert(db *gorm.DB, w *models.Weather) {

	db.Model(&w).Association("City").Find(&w.City)

	if w.Temperature >= 30 {
		alert := models.Alert{
			CityID:  w.City.ID,
			Message: fmt.Sprintf("Alerta de alta temperatura (%.0f°C) para la ciudad de %s", math.Ceil(w.Temperature), w.City.Name),
			Type:    "high-temp",
		}

		db.Create(&alert)
	}

	if w.Temperature <= 15 {
		alert := models.Alert{
			CityID:  w.City.ID,
			Message: fmt.Sprintf("Alerta de baja temperatura (%.0f°C) para la ciudad de %s", math.Ceil(w.Temperature), w.City.Name),
			Type:    "low-temp",
		}

		db.Create(&alert)
	}

	if w.Pressure <= 1010 {
		alert := models.Alert{
			CityID:  w.City.ID,
			Message: fmt.Sprintf("Alerta de sistema de baja presión (%d hPa) para la ciudad de %s", w.Pressure, w.City.Name),
			Type:    "low-press",
		}

		db.Create(&alert)
	}

	if w.Humidity >= 60 {
		alert := models.Alert{
			CityID:  w.City.ID,
			Message: fmt.Sprintf("Detectados altos niveles de humedad (%d%%) para la ciudad de %s", w.Humidity, w.City.Name),
			Type:    "high-hum",
		}

		db.Create(&alert)
	}
}

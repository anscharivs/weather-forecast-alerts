package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	config "github.com/anscharivs/weather-forecast-alerts/internal"
	"github.com/anscharivs/weather-forecast-alerts/models"
	"gorm.io/gorm"
)

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
		url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=es", city.Name, cfg.APIKey)

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

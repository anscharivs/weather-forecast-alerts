package weather

import (
	"encoding/json"
	"fmt"
	"net/http"

	config "github.com/anscharivs/weather-forecast-alerts/internal"
	"github.com/anscharivs/weather-forecast-alerts/models"
	"gorm.io/gorm"
)

type apiResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func FetchAndStoreWeatherData(db *gorm.DB, cfg config.Config) {

	var cities []models.City

	db.Find(&cities)

	for _, city := range cities {
		url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city.Name, cfg.APIKey)

		res, err := http.Get(url)

		if err != nil {
			fmt.Println("Error for ", city.Name, " weather ", err)
			continue
		}

		defer res.Body.Close()

		var data apiResponse

		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			fmt.Println("Response decoding error", err)
			continue
		}

		weather := models.Weather{
			CityID:      city.ID,
			Temperature: data.Main.Temp,
			Humidity:    data.Main.Humidity,
			Description: data.Weather[0].Description,
		}

		db.Create(&weather)

		fmt.Printf("Saved %s: %.1f°C, %s\n", city.Name, weather.Temperature, weather.Description)
	}

}

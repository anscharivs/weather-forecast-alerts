package api

import (
	"errors"
	"fmt"
	"math"
	"net/http"

	config "github.com/anscharivs/weather-forecast-alerts/internal"
	"github.com/anscharivs/weather-forecast-alerts/internal/weather"
	"github.com/anscharivs/weather-forecast-alerts/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/goodsign/monday"
	"gorm.io/gorm"
)

type WeatherView struct {
	CityName                string
	TemperatureInCelsius    float64
	MinTemperatureInCelsius float64
	MaxTemperatureInCelsius float64
	PressureInhPa           int
	HumidityInPercentage    int
	VisibilityInKm          int
	Description             string
	IconURL                 string
	FetchedAt               string
}

type AlertView struct {
	CityName  string
	CreatedAt string
	Message   string
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {

	r.GET("/", func(c *gin.Context) {

		var weathers []models.Weather

		subQuery := db.Model(&models.Weather{}).
			Select("MAX(id)").
			Group("city_id")

		if err := db.Preload("City").
			Where("id IN (?)", subQuery).
			Find(&weathers).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var views []WeatherView

		for _, w := range weathers {
			views = append(views, WeatherView{
				CityName:                w.City.Name,
				TemperatureInCelsius:    math.Ceil(w.Temperature),
				MinTemperatureInCelsius: math.Ceil(w.MinTemperature),
				MaxTemperatureInCelsius: math.Ceil(w.MaxTemperature),
				PressureInhPa:           w.Pressure,
				HumidityInPercentage:    w.Humidity,
				VisibilityInKm:          w.Visibility / 1000,
				Description:             w.Description,
				IconURL:                 fmt.Sprintf("https://openweathermap.org/img/wn/%s@2x.png", w.Icon),
				FetchedAt:               monday.Format(w.FetchedAt, "Monday 02 January 2006", monday.LocaleEsES),
			})
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"weathers": views,
		})
	})

	r.GET("/cities", func(c *gin.Context) {
		var cities []models.City
		db.Find(&cities)              // Search for records
		c.JSON(http.StatusOK, cities) // Serializes
	})

	r.GET("/weather", func(c *gin.Context) {
		// With query parameter  https://gin-gonic.com/en/docs/examples/querystring-param/#_top
		cityName := c.Query("city")

		if cityName != "" {
			var city models.City
			db.Where("name = ?", cityName).Find(&city)

			var weather models.Weather
			db.Where("city_id = ?", city.ID).Find(&weather)
			c.JSON(http.StatusOK, weather)
		} else {
			var weather []models.Weather
			db.Find(&weather)
			c.JSON(http.StatusOK, weather)
		}
	})

	r.GET("/fetch", func(c *gin.Context) {
		cfg := config.LoadConfig()
		weather.FetchAndStoreWeatherData(db, cfg)
		c.JSON(http.StatusOK, gin.H{"message": "Fetched data for cities"})
	})

	r.GET("/check-cities", func(c *gin.Context) {

		var cities []models.City
		var citiesCount int64
		db.Find(&cities).Count(&citiesCount)

		if citiesCount == 0 {

			city := models.City{Name: "Morelia"}
			db.Create(&city)
			cfg := config.LoadConfig()
			weather.FetchAndStoreWeatherData(db, cfg)

			c.JSON(http.StatusOK, gin.H{"no_records": true})
		} else {

			var weathers []models.Weather
			var weathersCount int64
			db.Find(&weathers).Count(&weathersCount)

			if weathersCount == 0 {

				cfg := config.LoadConfig()
				weather.FetchAndStoreWeatherData(db, cfg)
				c.JSON(http.StatusOK, gin.H{"no_records": true})

			} else {
				c.JSON(http.StatusOK, gin.H{"no_records": false})
			}
		}
	})

	r.GET("/registers", func(c *gin.Context) {

		var weathers []models.Weather

		db.Preload("City").
			Order("fetched_at DESC").
			Find(&weathers)

		var views []WeatherView

		for _, w := range weathers {
			views = append(views, WeatherView{
				CityName:                w.City.Name,
				TemperatureInCelsius:    math.Ceil(w.Temperature),
				MinTemperatureInCelsius: math.Ceil(w.MinTemperature),
				MaxTemperatureInCelsius: math.Ceil(w.MaxTemperature),
				PressureInhPa:           w.Pressure,
				HumidityInPercentage:    w.Humidity,
				VisibilityInKm:          w.Visibility / 1000,
				Description:             w.Description,
				IconURL:                 fmt.Sprintf("https://openweathermap.org/img/wn/%s@2x.png", w.Icon),
				FetchedAt:               monday.Format(w.FetchedAt, "02/01/2006 15:04", monday.LocaleEsES),
			})
		}

		c.HTML(http.StatusOK, "registers.html", gin.H{
			"weathers": views,
		})
	})

	r.GET("/alerts", func(c *gin.Context) {

		var alerts []models.Alert

		subQuery := db.Model(&models.Alert{}).
			Select("MAX(id)").
			Group("city_id, type")

		if err := db.Preload("City").
			Where("id IN (?)", subQuery).
			Order("created_at DESC").
			Find(&alerts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var views []AlertView

		for _, w := range alerts {
			views = append(views, AlertView{
				CityName:  w.City.Name,
				CreatedAt: monday.Format(w.CreatedAt, "02/01/2006 15:04", monday.LocaleEsES),
				Message:   w.Message,
			})
		}

		c.HTML(http.StatusOK, "alerts.html", gin.H{
			"alerts": views,
		})
	})

	r.GET("/new-city", func(c *gin.Context) {
		c.HTML(http.StatusOK, "new-city.html", gin.H{})
	})

	r.POST("/new-city", func(c *gin.Context) {

		type NewRegister struct {
			Name string `form:"name" binding:"required,max=50"`
		}

		var form NewRegister

		if err := c.ShouldBind(&form); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				errores := make(map[string]string)
				for _, fe := range ve {
					field := fe.Field()
					tag := fe.Tag()

					switch tag {
					case "required":
						errores[field] = "Este campo es obligatorio"
					case "max":
						errores[field] = "Debe tener como máximo " + fe.Param() + " caracteres"
					default:
						errores[field] = "Valor inválido"
					}
				}
				var erroresStr string

				for _, msg := range errores {
					if erroresStr != "" {
						erroresStr += ", "
					}
					erroresStr += msg
				}

				c.JSON(http.StatusBadRequest, gin.H{"error": erroresStr})

				return
			}

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if weather.ExistsInDB(db, form.Name) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "El nombre ya está registrado"})
			return
		}

		city := models.City{Name: form.Name}

		db.Create(&city)
		cfg := config.LoadConfig()
		weather.FetchAndStoreWeatherData(db, cfg)

		c.JSON(http.StatusOK, gin.H{
			"name": form.Name,
		})
	})

	r.GET("/delete-weather", func(c *gin.Context) {
		db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Weather{})
		c.JSON(http.StatusOK, gin.H{"message": "All weather records deleted"})
	})
}

package api

import (
	"net/http"

	"github.com/anscharivs/weather-forecast-alerts/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/cities", func(c *gin.Context) {
		var cities []models.City
		db.Find(&cities)              // Search for records
		c.JSON(http.StatusOK, cities) // Serializes
	})

	r.GET("/weather", func(c *gin.Context) {
		// With query parameter  https://gin-gonic.com/en/docs/examples/querystring-param/#_top
		cityName := c.Query("city")
		var city models.City
		db.Where("name = ?", cityName).Find(&city)

		var weather models.Weather
		db.Where("city_id = ?", city.ID).Find(&weather)
		c.JSON(http.StatusOK, weather)
	})
}

package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {

	r.GET("/", IndexHandler(db))
	r.GET("/check-cities", CheckCitiesHandler(db))
	r.GET("/registers", RegistersHandler(db))
	r.GET("/alerts", AlertsHandler(db))
	r.GET("/new-city", NewCityGetHandler())
	r.POST("/new-city", NewCityPostHandler(db))
	r.DELETE("/delete-city", DeleteCityHandler(db))

	// Manual
	r.GET("/cities", CitiesHandler(db))
	r.GET("/weather", WeatherHandler(db))
	r.GET("/fetch", FetchHandler(db))
	r.GET("/delete-weather", DeleteAllWeatherHandler(db))
}

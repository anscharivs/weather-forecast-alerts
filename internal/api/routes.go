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
}

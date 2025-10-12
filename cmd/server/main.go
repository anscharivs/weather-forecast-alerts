package main

import (
	"time"

	"github.com/anscharivs/weather-forecast-alerts/database"
	config "github.com/anscharivs/weather-forecast-alerts/internal"
	"github.com/anscharivs/weather-forecast-alerts/internal/api"
	"github.com/anscharivs/weather-forecast-alerts/internal/weather"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.InitDB(cfg)

	if err != nil {
		return
	}

	weather.StartWeatherPolling(db, cfg, 15*time.Minute) // Goroutine

	r := gin.Default()

	r.SetTrustedProxies(nil)

	r.LoadHTMLGlob("templates/*")

	api.RegisterRoutes(r, db)

	r.Run(":8080") // Run local server
}

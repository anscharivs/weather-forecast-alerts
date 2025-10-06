package city

import (
	"fmt"

	"github.com/anscharivs/weather-forecast-alerts/database"
	config "github.com/anscharivs/weather-forecast-alerts/internal"
	"github.com/anscharivs/weather-forecast-alerts/internal/weather"
	"github.com/spf13/cobra"
)

var FetchWeatherCmd = &cobra.Command{
	Use:   "fetch-weather",
	Short: "Fetch current weather data for all cities",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		db, err := database.InitDB(cfg)
		if err != nil {
			fmt.Println("Error with DB", err)
			return
		}

		weather.FetchAndStoreWeatherData(db, cfg)
	},
}

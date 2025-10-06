package city

import (
	"fmt"

	"github.com/anscharivs/weather-forecast-alerts/database"
	config "github.com/anscharivs/weather-forecast-alerts/internal"
	"github.com/anscharivs/weather-forecast-alerts/models"
	"github.com/spf13/cobra"
)

var ListCitiesCmd = &cobra.Command{
	Use:   "list-cities",
	Short: "List of registered cities",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.LoadConfig()

		db, err := database.InitDB(cfg)

		if err != nil {
			fmt.Println("Error with DB", err)
			return
		}

		var cities []models.City

		db.Find(&cities)

		for _, city := range cities {
			fmt.Println("-", city, city.Name)
		}
	},
}

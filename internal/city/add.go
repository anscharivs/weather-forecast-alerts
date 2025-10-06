package city

import (
	"fmt"

	"github.com/anscharivs/weather-forecast-alerts/database"
	config "github.com/anscharivs/weather-forecast-alerts/internal"
	"github.com/anscharivs/weather-forecast-alerts/models"
	"github.com/spf13/cobra"
)

var AddCityCmd = &cobra.Command{
	Use:   "add-city",
	Short: "Add a city",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You should put name city")
		}

		cfg := config.LoadConfig()

		db, err := database.InitDB(cfg)

		if err != nil {
			fmt.Println("Error with DB", err)
			return
		}

		city := models.City{Name: args[0]}

		db.Create(&city)

		fmt.Println(city.Name, " Added")
	},
}

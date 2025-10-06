package main

import (
	"fmt"
	"os"

	"github.com/anscharivs/weather-forecast-alerts/internal/city"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "weather-alerts"}

	// Command register
	rootCmd.AddCommand(city.AddCityCmd)
	rootCmd.AddCommand(city.ListCitiesCmd)
	rootCmd.AddCommand(city.FetchWeatherCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

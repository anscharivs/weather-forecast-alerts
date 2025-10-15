package weather_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	config "github.com/anscharivs/weather-forecast-alerts/internal"
	"github.com/anscharivs/weather-forecast-alerts/internal/weather"
	"github.com/anscharivs/weather-forecast-alerts/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *gorm.DB {

	dialector := sqlite.Dialector{
		DSN:        "file::memory:?cache=shared",
		DriverName: "sqlite",
	}

	db, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}

	db.AutoMigrate(&models.City{}, &models.Weather{}, &models.Alert{})

	return db
}

func TestFetchAndStoreWeatherData(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"main":{"temp":36.5,"temp_min":36.5,"temp_max":36.5,"humidity":80,"pressure":1010},"weather":[{"description":"cielo claro","icon":"01n"}],"visibility":10000,"dt":1760496964}`))
	}))

	defer mockServer.Close()

	cfg := config.Config{
		APIKey: "tuza_123",
	}

	db := setupTestDB(t)

	db.Create(&models.City{Name: "Ciudad Tuza"})

	weather.BaseURL = mockServer.URL

	weather.FetchAndStoreWeatherData(db, cfg)

	var weathers []models.Weather

	db.Find(&weathers)

	if len(weathers) != 1 {
		t.Errorf("Expected 1 weather record, got %d", len(weathers))
	}

	var alerts []models.Alert

	db.Find(&alerts)

	if len(alerts) != 3 {
		t.Errorf("Expected 3 alerts, got %d", len(alerts))
	}
}

package configs

import (
	"encoding/json"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
)

const (
	DefaultMaxRides             = 100
	DefaultLocationUpdateRateMs = 5000
)

type Config struct {
	RideRequestRateMs    int `json:"rideRequestRateMs" validate:"min=0"`
	MaxRides             int `json:"maxRides" validate:"min=0"`
	DriverCnt            int `json:"driverCnt" validate:"min=0"` // if 0, spawn a driver next to each rider
	CancelPct            int `json:"cancelPct" validate:"min=0"` // percent chance that a ride will be canceled each second while an accept/reject decision has yet to be made
	RejectPct            int `json:"rejectPct" validate:"min=0"` // percent chance that a ride request is rejected
	DecisionSpeedMs      int `json:"decisionSpeedMs" validate:"min=0"`
	LocationUpdateRateMs int `json:"locationUpdateRateMs" validate:"min=1000"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func GetConfig(path string) (config Config) {
	if file, err := os.Open(path); err != nil {
		log.Printf("%s; failed to load file at path %s; falling back to default\n", err, path)
	} else if err := json.NewDecoder(file).Decode(&config); err != nil {
		log.Printf(" %s; failed to load file at path %s; falling back to default\n", err, path)
	} else {
		defer file.Close()
		if err := validate.Struct(config); err != nil {
			log.Printf("%s; invalid file at path %s; falling back to default\n", err, path)
			config = Config{}
		}
	}

	if config.MaxRides == 0 {
		config.MaxRides = DefaultMaxRides
	}
	if config.LocationUpdateRateMs == 0 {
		config.LocationUpdateRateMs = DefaultLocationUpdateRateMs
	}

	log.Printf("%#v\n", config)
	return config
}

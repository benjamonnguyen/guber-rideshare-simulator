package configs

import (
	"encoding/json"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
)

var currConfig config

const (
	defaultWebSocketServerAddr  = "localhost:8081"
	defaultMaxRides             = 100
	defaultLocationUpdateRateMs = 5000
)

type config struct {
	WebSocketServerAddr  string `json:"webSocketServerAddr"`
	RideRequestRateMs    int    `json:"rideRequestRateMs" validate:"min=0"`
	MaxRides             int    `json:"maxRides" validate:"min=0"`
	DriverCnt            int    `json:"driverCnt" validate:"min=0"` // if 0, spawn a driver next to each rider
	CancelPct            int    `json:"cancelPct" validate:"min=0"` // percent chance that a ride will be canceled each second while an accept/reject decision has yet to be made
	RejectPct            int    `json:"rejectPct" validate:"min=0"` // percent chance that a ride request is rejected
	DecisionSpeedMs      int    `json:"decisionSpeedMs" validate:"min=0"`
	LocationUpdateRateMs int    `json:"locationUpdateRateMs" validate:"min=1000"`
}

var validate = validator.New(validator.WithRequiredStructEnabled())

func GetConfig() config {
	return currConfig
}

func LoadConfig(path string) config {
	if file, err := os.Open(path); err != nil {
		log.Printf("%s; failed to load file at path %s; falling back to default\n", err, path)
	} else if err := json.NewDecoder(file).Decode(&currConfig); err != nil {
		log.Printf(" %s; failed to load file at path %s; falling back to default\n", err, path)
	} else {
		defer file.Close()
		if err := validate.Struct(currConfig); err != nil {
			log.Printf("%s; invalid file at path %s; falling back to default\n", err, path)
			currConfig = config{}
		}
	}

	if currConfig.WebSocketServerAddr == "" {
		currConfig.WebSocketServerAddr = defaultWebSocketServerAddr
	}
	if currConfig.MaxRides == 0 {
		currConfig.MaxRides = defaultMaxRides
	}
	if currConfig.LocationUpdateRateMs == 0 {
		currConfig.LocationUpdateRateMs = defaultLocationUpdateRateMs
	}

	log.Printf("%#v\n", currConfig)
	return currConfig
}

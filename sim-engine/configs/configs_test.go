package configs

import (
	"testing"
)

func TestGetConfigWithNoConfigPath(t *testing.T) {
	conf := GetConfig("")

	if conf.MaxRides != DefaultMaxRides {
		t.Errorf("maxRides incorrect, got: %d, want: %d", conf.MaxRides, DefaultMaxRides)
	}
	if conf.LocationUpdateRateMs != DefaultLocationUpdateRateMs {
		t.Errorf("locationUpdateRateMs incorrect, got: %d, want: %d", conf.LocationUpdateRateMs, DefaultLocationUpdateRateMs)
	}
}

func TestGetConfigWithExampleConfigPath(t *testing.T) {
	conf := GetConfig("example-config.json")

	got := conf.RideRequestRateMs
	want := 5000
	if got != want {
		t.Errorf("RideRequestRateMs incorrect, got: %d, want: %d", got, want)
	}
	got = conf.MaxRides
	want = 100
	if got != want {
		t.Errorf("MaxRides incorrect, got: %d, want: %d", got, want)
	}
	got = conf.DriverCnt
	want = 50
	if got != want {
		t.Errorf("DriverCnt incorrect, got: %d, want: %d", got, want)
	}
	got = conf.CancelPct
	want = 10
	if got != want {
		t.Errorf("CancelPct incorrect, got: %d, want: %d", got, want)
	}
	got = conf.RejectPct
	want = 20
	if got != want {
		t.Errorf("RejectPct incorrect, got: %d, want: %d", got, want)
	}
	got = conf.DecisionSpeedMs
	want = 15000
	if got != want {
		t.Errorf("DecisionSpeedMs incorrect, got: %d, want: %d", got, want)
	}
	got = conf.LocationUpdateRateMs
	want = 5000
	if got != want {
		t.Errorf("LocationUpdateRateMs incorrect, got: %d, want: %d", got, want)
	}
}

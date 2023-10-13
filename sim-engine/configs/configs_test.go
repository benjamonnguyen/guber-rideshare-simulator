package configs

import (
	"testing"
)

func TestLoadConfigWithNoConfigPath(t *testing.T) {
	conf := LoadConfig("")

	if conf.WebSocketServerAddr != defaultWebSocketServerAddr {
		t.Errorf("webSocketServerAddr incorrect, got: %s, want: %s", conf.WebSocketServerAddr, defaultWebSocketServerAddr)
	}
	if conf.MaxRides != defaultMaxRides {
		t.Errorf("maxRides incorrect, got: %d, want: %d", conf.MaxRides, defaultMaxRides)
	}
	if conf.LocationUpdateRateMs != defaultLocationUpdateRateMs {
		t.Errorf("locationUpdateRateMs incorrect, got: %d, want: %d", conf.LocationUpdateRateMs, defaultLocationUpdateRateMs)
	}
}

func TestLoadConfigWithExampleConfigPath(t *testing.T) {
	conf := LoadConfig("example-config.json")

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

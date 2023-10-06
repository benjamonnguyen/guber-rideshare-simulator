package ridemodel

type Ride struct {
	ID          string   `json:"id"`
	PassengerID string   `json:"passengerId"`
	DriverID    string   `json:"driverId,omitempty"`
	Status      Status   `json:"status"`
	Pickup      [2]int   `json:"pickup"`
	Dropoff     [2]int   `json:"dropoff"`
	PickupPath  [][2]int `json:"pickupPath,omitempty"`
	DropoffPath [][2]int `json:"dropoffPath,omitempty"`
}

type RideRequest struct {
	PassengerID string `json:"passengerId"`
	Pickup      [2]int `json:"pickup"`
	Dropoff     [2]int `json:"dropoff"`
}

type Status int

const (
	StatusPending Status = iota
	StatusEnrouteToPickup
	StatusEnrouteToDropoff
	StatusCancelled
	StatusComplete
)

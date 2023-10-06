package ridemodel

import "github.com/benjamonnguyen/guber-ridershare-simulator/model/coord"

type Ride struct {
	ID          string        `json:"id"`
	PassengerID string        `json:"passengerId"`
	DriverID    string        `json:"driverId,omitempty"`
	Status      Status        `json:"status"`
	Pickup      coord.Coord   `json:"pickup"`
	Dropoff     coord.Coord   `json:"dropoff"`
	PickupPath  []coord.Coord `json:"pickupPath,omitempty"`
	DropoffPath []coord.Coord `json:"dropoffPath,omitempty"`
}

type RideRequest struct {
	PassengerID string      `json:"passengerId"`
	Pickup      coord.Coord `json:"pickup"`
	Dropoff     coord.Coord `json:"dropoff"`
}

type Status int

const (
	StatusPending Status = iota
	StatusEnrouteToPickup
	StatusEnrouteToDropoff
	StatusCancelled
	StatusComplete
)

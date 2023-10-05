package driver

type Driver struct {
	ID     string       `json:"id"`
	Name   string       `json:"name"`
	Status DriverStatus `json:"status"`
}

type DriverStatus int

const (
	OFFLINE DriverStatus = iota
	ONLINE
	ACTIVE
	RIDE_IN_PROGRESS
)

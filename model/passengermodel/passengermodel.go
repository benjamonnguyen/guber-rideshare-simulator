package passengermodel

type Passenger struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status Status `json:"status"`
}

type Status int

const (
	StatusInactive Status = iota
	StatusRequesting
	StatusEnroute
)

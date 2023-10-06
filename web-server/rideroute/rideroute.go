package rideroute

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/benjamonnguyen/guber-ridershare-simulator/model/coord"
	"github.com/benjamonnguyen/guber-ridershare-simulator/model/drivermodel"
	"github.com/benjamonnguyen/guber-ridershare-simulator/model/passengermodel"
	"github.com/benjamonnguyen/guber-ridershare-simulator/model/ridemodel"
	"github.com/benjamonnguyen/guber-ridershare-simulator/web-server/driverroute"
	"github.com/benjamonnguyen/guber-ridershare-simulator/web-server/passengerroute"
	"github.com/julienschmidt/httprouter"
)

var dummyCollection = map[string]*ridemodel.Ride{
	"1": {
		ID:          "1",
		PassengerID: "1",
		DriverID:    "1",
		Status:      ridemodel.StatusComplete,
		Pickup:      coord.Coord{X: 0, Y: 0},
		Dropoff:     coord.Coord{X: 0, Y: 0},
		PickupPath:  []coord.Coord{{X: 0, Y: 0}},
		DropoffPath: []coord.Coord{{X: 0, Y: 0}},
	},
}

func GetRide(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ride, ok := dummyCollection[ps.ByName("id")]
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	json, err := json.Marshal(ride)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(json)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func CreateRide(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var rideRequest ridemodel.RideRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rideRequest); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	passenger, ok := passengerroute.DummyCollection[rideRequest.PassengerID]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	id := fmt.Sprintf("%d", len(dummyCollection)+1)

	ride := ridemodel.Ride{
		ID:          id,
		PassengerID: rideRequest.PassengerID,
		Pickup:      rideRequest.Pickup,
		Dropoff:     rideRequest.Dropoff,
		Status:      ridemodel.StatusPending,
	}

	// TODO transactional
	dummyCollection[id] = &ride
	passenger.Status = passengermodel.StatusRequesting

	// TODO async call to Dispatcher

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(id)); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func Accept(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	passenger, ride, statusCode, err := validateUpdateRequest(r, ps)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	driverId, statusCode, err := getDriverIdFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	driver, ok := driverroute.DummyCollection[driverId]
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	// TODO call to Dispatcher to accept and acquire lock { rideId, driverId } - last opportunity to check for cancel

	// TODO transactional
	passenger.Status = passengermodel.StatusEnroute
	driver.Status = drivermodel.StatusEnroute
	ride.Status = ridemodel.StatusEnrouteToPickup
	ride.DriverID = driver.ID

	// TODO async call to Dispatcher
}

func Reject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, _, statusCode, err := validateUpdateRequest(r, ps)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	driverId, statusCode, err := getDriverIdFromBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	driver, ok := driverroute.DummyCollection[driverId]
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	driver.Status = drivermodel.StatusActive

	// async call to Dispatcher
}

func Cancel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	passenger, ride, statusCode, err := validateUpdateRequest(r, ps)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	if ride.Status != ridemodel.StatusPending {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// TODO transactional
	ride.Status = ridemodel.StatusCancelled
	passenger.Status = passengermodel.StatusInactive

	// TODO async call to Dispatcher to notify and update Driver
}

func Complete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	passenger, ride, statusCode, err := validateUpdateRequest(r, ps)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	driver, ok := driverroute.DummyCollection[ride.DriverID]
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	// TODO transactional
	ride.Status = ridemodel.StatusComplete
	passenger.Status = passengermodel.StatusInactive
	driver.Status = drivermodel.StatusActive

	// TODO async call to Dispatcher to credit/debit driver/passenger
}

func validateUpdateRequest(r *http.Request, ps httprouter.Params) (passenger *passengermodel.Passenger, ride *ridemodel.Ride, statusCode int, err error) {
	// TODO should these reads acquire a lock for concurrency safety?
	ride, ok := dummyCollection[ps.ByName("id")]
	if !ok {
		return nil, nil, http.StatusNotFound, errors.New("")
	}
	if ride.Status != ridemodel.StatusPending {
		return nil, nil, http.StatusBadRequest, errors.New("")
	}

	passenger, ok = passengerroute.DummyCollection[ride.PassengerID]
	if !ok {
		return nil, nil, http.StatusNotFound, errors.New("")
	}

	return passenger, ride, http.StatusOK, nil
}

func getDriverIdFromBody(body io.ReadCloser) (driverId string, statusCode int, err error) {
	driverIdData, err := io.ReadAll(body)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("")
	}
	if len(driverIdData) == 0 {
		errMsg := "no driverId provided"
		return "", http.StatusBadRequest, errors.New(errMsg)
	}

	return string(driverIdData), http.StatusOK, nil
}

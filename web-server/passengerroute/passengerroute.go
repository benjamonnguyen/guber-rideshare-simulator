package passengerroute

import (
	"encoding/json"
	"net/http"

	"github.com/benjamonnguyen/guber-ridershare-simulator/model/passengermodel"
	"github.com/julienschmidt/httprouter"
)

var DummyCollection = map[string]*passengermodel.Passenger{
	"1": {
		ID:   "1",
		Name: "Jesse",
	},
}

func GetPassenger(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	passenger, ok := DummyCollection[ps.ByName("id")]
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	json, err := json.Marshal(passenger)
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

func StreamLocation(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		// TODO StreamLocation
		w.Write([]byte(req.URL.String()))
	}
}

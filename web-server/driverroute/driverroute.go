package driverroute

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/benjamonnguyen/guber-ridershare-simulator/model/drivermodel"
	"github.com/julienschmidt/httprouter"
)

var DummyCollection = map[string]*drivermodel.Driver{
	"1": {
		ID:     "1",
		Name:   "Mike",
		Status: drivermodel.StatusInactive,
	},
}

// TODO dummyLocationStream
func dummyLocationStream(ctx context.Context) <-chan [2]int {
	stream := make(chan [2]int, 5)

	go func(ch chan<- [2]int) {
		defer close(stream)
		t := time.NewTicker(time.Second * 5)
		defer t.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				ch <- [2]int{0, 0}
			}
		}
	}(stream)

	return stream
}

func GetDriver(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	driver, ok := DummyCollection[ps.ByName("id")]
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	json, err := json.Marshal(*driver)
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

func StreamLocation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	stream := dummyLocationStream(r.Context())
	for coord := range stream {
		json, err := json.Marshal(coord)
		fmt.Println(string(json))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		_, err = w.Write(json)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		flusher.Flush()
	}
}

func EndSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	d, ok := DummyCollection[ps.ByName("id")]
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	d.Status = drivermodel.StatusInactive

	json, err := json.Marshal(d)
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

func StartSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	d, ok := DummyCollection[ps.ByName("id")]
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	d.Status = drivermodel.StatusActive

	json, err := json.Marshal(d)
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

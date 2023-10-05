package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/benjamonnguyen/guber-ridershare-simulator/web-server/driverroute"
	"github.com/benjamonnguyen/guber-ridershare-simulator/web-server/passengerroute"
	"github.com/julienschmidt/httprouter"
)

func main() {
	var addr = flag.String("addr", "localhost:8080", "http service address")
	flag.Parse()
	log.Println("starting web-server on", *addr)

	router := httprouter.New()
	// Driver routes
	router.GET("/driver/:id", driverroute.GetDriver)
	router.GET("/driver/:id/location", driverroute.StreamLocation)
	router.PUT("/driver/:id/start", driverroute.StartSession)
	router.PUT("/driver/:id/end", driverroute.EndSession)

	// Passenger routes
	router.GET("/passenger/:id", passengerroute.GetPassenger)
	router.GET("/passenger/:id/location", driverroute.StreamLocation) // TODO using same dummy stream

	log.Fatal(http.ListenAndServe(*addr, router))
}

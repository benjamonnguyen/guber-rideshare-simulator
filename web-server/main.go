package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/benjamonnguyen/guber-ridershare-simulator/web-server/driverroute"
	"github.com/benjamonnguyen/guber-ridershare-simulator/web-server/passengerroute"
	"github.com/benjamonnguyen/guber-ridershare-simulator/web-server/rideroute"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

func main() {
	const appName = "web-server"

	addr := flag.String("addr", "localhost:8080", "http service address")
	flag.Parse()

	router := httprouter.New()
	// Driver routes
	router.GET("/driver/:id", driverroute.GetDriver)
	router.GET("/driver/:id/location", driverroute.StreamLocation)
	router.PUT("/driver/:id/start", driverroute.StartSession)
	router.PUT("/driver/:id/end", driverroute.EndSession)

	// Passenger routes
	router.GET("/passenger/:id", passengerroute.GetPassenger)
	router.GET("/passenger/:id/location", driverroute.StreamLocation) // TODO using same dummy stream

	// Ride routes
	router.GET("/ride/:id", rideroute.GetRide)
	router.POST("/ride", rideroute.CreateRide)
	router.PUT("/ride/:id/accept", rideroute.Accept)
	router.PUT("/ride/:id/reject", rideroute.Reject)
	router.PUT("/ride/:id/cancel", rideroute.Cancel)
	router.PUT("/ride/:id/complete", rideroute.Complete)

	n := negroni.Classic()
	n.UseHandler(router)

	log.SetPrefix(fmt.Sprintf("[%s] ", appName))
	log.Printf("started %s on %s\n", appName, *addr)
	log.Fatal(http.ListenAndServe(*addr, n))
}

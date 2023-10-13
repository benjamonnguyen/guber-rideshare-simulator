// simulates client connections to the web-socket-server
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"sync"
	"time"

	"github.com/benjamonnguyen/guber-rideshare-simulator/sim-engine/configs"
	"github.com/benjamonnguyen/guber-rideshare-simulator/sim-engine/sim"
)

const appName = "sim-engine"

var clientsWg sync.WaitGroup

func main() {
	configPath := flag.String("conf", "", "path to config json file")
	flag.Parse()
	config := configs.LoadConfig(path.Join("configs", *configPath))

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt)
	locationUpdateTick := time.NewTicker(time.Millisecond * time.Duration(config.LocationUpdateRateMs)).C

	// TODO fetch and iterate through inactive drivers from DB
	for i := 1; i < config.DriverCnt+1; i++ {
		clientsWg.Add(1)
		go sim.SpawnDriver(fmt.Sprint(i), &clientsWg, locationUpdateTick, interruptSignal)
	}

	// TODO Timer to create ride requests

	log.SetPrefix(fmt.Sprintf("[%s] ", appName))

	clientsWg.Wait()
}

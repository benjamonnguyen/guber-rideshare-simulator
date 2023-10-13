package sim

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/benjamonnguyen/guber-ridershare-simulator/model/msgmodel"
	"github.com/benjamonnguyen/guber-rideshare-simulator/sim-engine/client"
	"github.com/benjamonnguyen/guber-rideshare-simulator/sim-engine/configs"
	"github.com/gorilla/websocket"
)

const writeTimeout = time.Second * 2
const readTimeout = time.Minute * 10

// SpawnDriver connects a client representing a Driver
func SpawnDriver(id string, clientsWg *sync.WaitGroup, locationUpdateTick <-chan time.Time, interruptSignal chan os.Signal) {
	defer clientsWg.Done()

	config := configs.GetConfig()

	var c client.Client
	c.Dial(fmt.Sprintf("ws://%s/connect/%s", config.WebSocketServerAddr, id), nil)
	defer c.Close()

	done := make(chan struct{})
	go handleServerMessages(&c, done)

	// write location updates
	for {
		select {
		case <-done:
			c.Shutdown(writeTimeout)
			return
		case <-interruptSignal:
			c.Shutdown(writeTimeout)
			return
		case <-locationUpdateTick:
			// TODO randomize coords
			locationMsg := msgmodel.LocationMsg{
				UserId: id,
				Coord:  [2]int{0, 0},
			}

			data, err := json.Marshal(locationMsg)
			if err != nil {
				log.Println("SpawnDriver: write:", err)
			}

			c.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := c.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("SpawnDriver: write:", err)
			}
		}
	}
}

func handleServerMessages(c *client.Client, done chan struct{}) {
	defer close(done)
	for {
		c.SetReadDeadline(time.Now().Add(readTimeout))
		_, msg, err := c.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("SpawnDriver: read: killing driver sim")
				return
			}
			log.Println("SpawnDriver: read:", err)
			// reconnect
			c.Close()
			c.Connect()
		} else {
			log.Println(string(msg))
		}
	}
}

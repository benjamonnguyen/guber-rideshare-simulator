package route

import (
	"log"
	"net/http"
	"time"

	"github.com/benjamonnguyen/guber-rideshare-simulator/web-socket-server/clientregistry"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

const writeTimeout = time.Second * 2
const readTimeout = time.Second * 10

var upgrader = websocket.Upgrader{}

var ticker *time.Ticker = time.NewTicker(time.Second * 30)

func Connect(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	clientId := ps.ByName("clientId")
	if len(clientId) == 0 {
		http.Error(w, "no clientId provided", http.StatusBadRequest)
		return
	}

	_, ok := clientregistry.ClientRegistry[clientId]
	if ok {
		http.Error(w, "", http.StatusConflict)
		return
	}

	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}

	clientregistry.ClientRegistry[clientId] = c

	done := make(chan struct{})

	// read
	go func() {
		for {
			c.SetReadDeadline(time.Now().Add(readTimeout))
			msgType, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Connect: read:", err)
				close(done)
				return
			}
			handleClientMsg(clientId, msgType, msg)
		}
	}()

	// heartbeat
	go func() {
		defer func() {
			if err := c.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(writeTimeout)); err != nil {
				c.Close()
			}
			delete(clientregistry.ClientRegistry, clientId)
			log.Printf("Connect: heartbeat: removed client %s from registry\n", clientId)
		}()
		failCnt := 0
		for {
			if failCnt >= 2 {
				log.Printf("Connect: heartbeat: pings to client %s failed\n", clientId)
				return
			}
			select {
			case <-done:
				return
			case <-ticker.C:
				if err := c.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeTimeout)); err != nil {
					failCnt++
				}
			}
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func handleClientMsg(clientId string, msgType int, msg []byte) {
	log.Printf("%s (clientId: %s, type: %d)", msg, clientId, msgType)
}

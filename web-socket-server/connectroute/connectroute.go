package connectroute

import (
	"log"
	"net/http"

	"github.com/benjamonnguyen/guber-rideshare-simulator/web-socket-server/clientregistry"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{}

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
		log.Println("upgrade:", err)
		return
	}

	clientregistry.ClientRegistry[clientId] = c

	go func() {
		defer c.Close()
		for {
			msgType, msg, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err) {
					delete(clientregistry.ClientRegistry, clientId)
					log.Printf("remove client %s from registry\n", clientId)
				}
				log.Println("read:", err)
				break
			}
			handleClientMsg(clientId, msgType, msg)
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func handleClientMsg(clientId string, msgType int, msg []byte) {
	log.Printf("%s (clientId: %s, type: %d)", msg, clientId, msgType)
}

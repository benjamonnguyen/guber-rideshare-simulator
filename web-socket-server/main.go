package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/benjamonnguyen/guber-ridershare-simulator/model/msgmodel"
	"github.com/benjamonnguyen/guber-rideshare-simulator/web-socket-server/clientregistry"
	"github.com/benjamonnguyen/guber-rideshare-simulator/web-socket-server/connectroute"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

func main() {
	var addr = flag.String("addr", "localhost:8081", "http service address")
	flag.Parse()

	router := httprouter.New()
	router.GET("/connect/:clientId", connectroute.Connect)
	router.POST("/msg", msgClient)

	n := negroni.Classic()
	n.UseHandler(router)

	log.Println("started web-socket-server on", *addr)
	log.Fatal(http.ListenAndServe(*addr, n))
}

func msgClient(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var serverMsg msgmodel.ServerMsg
	if err := json.NewDecoder(req.Body).Decode(&serverMsg); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	c, ok := clientregistry.ClientRegistry[serverMsg.ClientId]
	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	if err := c.WriteMessage(websocket.TextMessage, []byte(serverMsg.Msg)); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		// TODO retry logic
	}

	w.WriteHeader(http.StatusOK)
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/benjamonnguyen/guber-ridershare-simulator/model/msgmodel"
	"github.com/benjamonnguyen/guber-rideshare-simulator/web-socket-server/clientregistry"
	"github.com/benjamonnguyen/guber-rideshare-simulator/web-socket-server/route"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

const appName = "web-socket-server"

var addr = flag.String("addr", "localhost:8081", "http service address")

func main() {
	go handleGraceShutdown()
	flag.Parse()

	router := httprouter.New()
	router.GET("/connect/:clientId", route.Connect)
	router.POST("/msg", msgClient)

	n := negroni.Classic()
	n.UseHandler(router)

	log.SetPrefix(fmt.Sprintf("[%s] ", appName))
	log.Printf("started %s on %s\n", appName, *addr)
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

func handleGraceShutdown() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("graceful shutdown started")
	done := make(chan struct{}, 1)
	var wg sync.WaitGroup
	for id, conn := range clientregistry.ClientRegistry {
		id := id
		conn := conn
		wg.Add(1)
		go func() {
			conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second*2))
			log.Println("sent closeMessage to client", id)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-time.After(time.Second * 5):
		log.Println("graceful shutdown timeout")
	case <-done:
		log.Println("graceful shutdown complete")
	}
	os.Exit(0)
}

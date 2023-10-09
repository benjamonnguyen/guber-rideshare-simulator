package clientregistry

import "github.com/gorilla/websocket"

var ClientRegistry = make(map[string]*websocket.Conn)

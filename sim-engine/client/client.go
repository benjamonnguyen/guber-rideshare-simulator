// Package client provides a reconnectable client based on gorilla/websocket
package client

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	url       string
	reqHeader http.Header
	*websocket.Conn
}

func (client *Client) Dial(url string, reqHeader http.Header) {
	url, err := parseURL(url)

	if err != nil {
		log.Fatalln("Dial:", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(url, reqHeader)
	if err != nil {
		// TODO handle retry with backoff and timeout
		return
	}
	log.Println("Dial: connected to", url)
	client.Conn = conn
	client.url = url

	// TODO go keepAlive() care for leak
}

func (client *Client) Connect() {
	client.Dial(client.url, client.reqHeader)
}

// Shutdown gracefully closes the connection by sending the websocket.CloseMessage.
// The writeWait param defines the duration before the deadline of the write operation is hit.
func (client *Client) Shutdown(writeWait time.Duration) {
	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	err := client.WriteControl(websocket.CloseMessage, msg, time.Now().Add(writeWait))
	if err != nil && err != websocket.ErrCloseSent {
		log.Printf("Shutdown: %v", err)
		client.Close()
	}
}

func parseURL(urlStr string) (string, error) {
	if urlStr == "" {
		return "", errors.New("url cannot be empty")
	}

	u, err := url.Parse(urlStr)

	if err != nil {
		return "", errors.New("url: " + err.Error())
	}

	if u.Scheme != "ws" && u.Scheme != "wss" {
		return "", errors.New("url: websocket uris must start with ws or wss scheme")
	}

	if u.User != nil {
		return "", errors.New("url: user name and password are not allowed in websocket URIs")
	}

	return urlStr, nil
}

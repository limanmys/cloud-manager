package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func main() {

	ws_url := "wss://127.0.0.1:8211/ws"
	header := make(http.Header)
	header.Add("machine_id", "1111")
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}

	conn, _, err := dialer.Dial(ws_url, header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			fmt.Printf("Received: %s\n", message)
			conn.WriteMessage(websocket.TextMessage, []byte("OK"))
		}

	}()
	<-done
}

// client.go
package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type WSClient struct {
	id          string
	out         chan []byte
	broadcaster chan []byte
	conn        *websocket.Conn
}

type packet struct {
	id   string
	data []byte
}

func NewWSClient(id string, broadcaster chan []byte, conn *websocket.Conn) *WSClient {
	c := new(WSClient)
	c.id = id
	c.out = make(chan []byte)
	c.broadcaster = broadcaster
	c.conn = conn
	return c
}

func (ws *WSClient) HandleWebsocketConnection() {
	daytime := time.Now().String() // get current time
	err := ws.conn.WriteMessage(websocket.BinaryMessage, []byte(daytime+"\n"))

	if err != nil {
		log.Println(err)
	}

	// goroutine for receiving data from client
	go func(ws *WSClient) {

		for {
			_, data, err := ws.conn.ReadMessage()

			if err != nil {
				log.Println(err)
				break
			}

			ws.broadcaster <- data
		}

	}(ws)

	// goroutine for sending messages to client from server
	go func(ws *WSClient) {

		for {
			select {
			case data := <-ws.out:
				err := ws.conn.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					log.Println(err)
					break
				}
			}
		}

	}(ws)

}

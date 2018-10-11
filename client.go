// client.go
package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Packet struct {
	id   string
	data []byte
}

type WSClient struct {
	id          string
	out         chan Packet
	broadcaster chan Packet
	leave       chan string
	conn        *websocket.Conn
}

func NewWSClient(id string, broadcaster chan Packet,
	leave chan string, conn *websocket.Conn) *WSClient {

	c := new(WSClient)
	c.id = id
	c.out = make(chan Packet)
	c.broadcaster = broadcaster
	c.leave = leave
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

		defer func() {
			ws.conn.Close()
			ws.leave <- ws.id
		}()

		for {
			_, data, err := ws.conn.ReadMessage()

			if err != nil {
				log.Println(err)
				break
			}

			ws.broadcaster <- Packet{id: ws.id, data: data}
		}

	}(ws)

	// goroutine for sending messages to client from server
	go func(ws *WSClient) {

		defer func() {
			ws.conn.Close()
			ws.leave <- ws.id
		}()

		for {
			select {
			case packet := <-ws.out:
				msg := append([]byte(packet.id), packet.data...)
				err := ws.conn.WriteMessage(websocket.BinaryMessage, msg)
				if err != nil {
					log.Println(err)
					break
				}
			}
		}

	}(ws)

}

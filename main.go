// httpserver project main.go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const SOCKET_ID_LENGTH = 10

var (
	idValues []string = strings.Split("abcdefghilmnopqrstuvzABCDEFGHILMNOPQRSTUVZ1234567890xykjXYKJ", "")
	upgraderChat          = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgraderFileStream = websocket.Upgrader{
		ReadBufferSize: 
	}
	broadcaster = make(chan Packet)
	leave       = make(chan string)
	clients     = NewConcurrentMap()
)


func shuffle(a []string) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func createSocketId(vals []string) string {

	rand.Seed(time.Now().UnixNano())
	shuffle(vals)
	id := make([]string, SOCKET_ID_LENGTH)
	for i := 0; i < SOCKET_ID_LENGTH; i++ {
		index := rand.Intn(len(vals) - 1)
		id[i] = vals[index]
	}
	return strings.Join(id, "")

}

func getPort() string {
	port := os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return port
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Go http server here! Your public IP is: %s", r.RemoteAddr)
}

func handleWSRequest(w http.ResponseWriter, r *http.Request) {
	// Upgrade http GET request to a websocket

	ws, err := upgraderChat.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	id := createSocketId(idValues)
	log.Println("New connection: " + id)
	client := NewWSClient(id, broadcaster, leave, ws)
	clients.Add(id, client)
	client.HandleWebsocketConnection()

}

/*func handleWSFileStream(w http.ResponseWriter, r *http.Request) {
	//Set up new connection for streaming files
	
	fsws, err := upgraderFileStream.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	
	
	
}*/

func MasterWebsocketHandler() {
	for {
		select {
		case clientToDelete := <-leave:
			clients.Delete(clientToDelete)
			log.Println("Client " + clientToDelete + " deleted")

			clients.Iterate(func(id string, client *WSClient) {
				client.conn.WriteMessage(websocket.BinaryMessage,
					[]byte(clientToDelete+" disconnected\n"))
			})

		case packet := <-broadcaster:
			clients.Iterate(func(id string, client *WSClient) {
				if id != packet.id {
					client.out <- packet
				}
			})
		}
	}
}

func main() {
	port := getPort()

	go MasterWebsocketHandler()

	http.HandleFunc("/", handler)
	http.HandleFunc("/ws", handleWSRequest)
	//http.HandleFunc("/fs/ws", handleWSFileStream)
	fmt.Println("HTTP Server listening on port: " + port + "...")
	err := http.ListenAndServe("0.0.0.0:"+port, nil)

	if err != nil {
		log.Fatalln(err)
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/bobcob7/tech-debt/messages"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

// WSClient is how a Websocket Client is accessed
type WSClient struct {
	connection *websocket.Conn
	receiver   chan []byte
}

// Send a message to this particular client
func (client WSClient) Send(message []byte) error {
	return client.connection.WriteMessage(websocket.BinaryMessage, message)
}

// Listen for messages from this client, and put them on the reciever channel
func (client WSClient) Listen() {
	defer client.connection.Close()
	for {
		mt, message, err := client.connection.ReadMessage()
		if mt != websocket.BinaryMessage {
			if err != nil {
				log.Println("Weird message error:", err)
				return
			}
		}
		if err != nil {
			log.Println("read error:", err)
			return
		}
		log.Printf("recv: %s", message)
		client.receiver <- message
	}
}

// Hub is where all clients are kept and broadcasts are handled
type Hub struct {
	clients  []WSClient
	Receiver chan []byte
}

// New creates a Hub
func New() Hub {
	return Hub{
		clients:  make([]WSClient, 10),
		Receiver: make(chan []byte, 10),
	}
}

// AddClient adds a new client to the broadcast group
func (hub *Hub) AddClient(connection *websocket.Conn) {
	newClient := WSClient{connection, hub.Receiver}
	go newClient.Listen()
	hub.clients = append(hub.clients, newClient)
}

// Broadcast a message to every client
func (hub *Hub) Broadcast(message []byte) []error {
	errors := make([]error, 0)
	for _, client := range hub.clients {
		err := client.Send(message)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

var hub = New()
var sum = 0

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	hub.AddClient(c)
}

// HubCounter recieves messages and broadcast their sums
func HubCounter(h Hub) {
	for {
		message := <-h.Receiver
		fmt.Println(message)
	}
}

var currentDebt = &messages.TechDebt{}

func checker(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: checker,
}

func upgrade(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	log.Printf("Got connection from %s", r.RemoteAddr)

	// Immediately send current tech debt
	var buffer []byte
	buffer, err = proto.Marshal(currentDebt)
	if err != nil {
		log.Println("marshal err:", err)
		return
	}
	err = c.WriteMessage(1, buffer)
	if err != nil {
		log.Println("write err:", err)
		return
	}

	var receiveMsg = &messages.TechDebt{}
	for {
		// Get new message
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read err:", err)
			break
		}
		err = proto.Unmarshal(message, receiveMsg)
		if err != nil {
			log.Println("unmarshal err:", err)
			break
		}

		// Add new tech debt
		currentDebt.TechDebt += receiveMsg.TechDebt
		fmt.Printf("%s added %d points:\tCurrent Score: %d\n", receiveMsg.Username, receiveMsg.TechDebt, currentDebt.TechDebt)

		// Write tech debt
		buffer, err := proto.Marshal(currentDebt)
		if err != nil {
			log.Println("marshal err:", err)
			break
		}
		err = c.WriteMessage(mt, buffer)
		if err != nil {
			log.Println("write err:", err)
			break
		}
	}
}

func main() {
	flag.Parse()
	currentDebt.TechDebt = 0
	log.SetFlags(0)
	http.HandleFunc("/ws", upgrade)
	log.Println("Upgrade using ws://127.0.0.1:3000/ws")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/bobcob7/tech-debt/messages"
	"github.com/golang-collections/collections/stack"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

// WSClient is how a Websocket Client is accessed
type WSClient struct {
	RemoteAddr string
	connection *websocket.Conn
	receiver   chan []byte
}

// Send a message to this particular client
func (client WSClient) Send(message []byte) error {
	if client.connection != nil {
		return client.connection.WriteMessage(websocket.BinaryMessage, message)
	}
	return errors.New("no client connection")
}

// Listen for messages from this client, and put them on the reciever channel
func (client WSClient) Listen() {
	//defer client.connection.Close()
	for {
		mt, message, err := client.connection.ReadMessage()
		if err != nil {
			log.Println("Weird message error:", err)
			return
		}
		if mt != websocket.BinaryMessage {
			if err != nil {
				log.Println("Weird message error 2:", err)
				return
			}
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
		clients:  make([]WSClient, 0, 10),
		Receiver: make(chan []byte, 10),
	}
}

// AddClient adds a new client to the broadcast group
func (hub *Hub) AddClient(remoteAddr string, connection *websocket.Conn) {
	log.Printf("Adding client for %s\n", remoteAddr)
	newClient := WSClient{remoteAddr, connection, hub.Receiver}
	go newClient.Listen()
	hub.clients = append(hub.clients, newClient)
	log.Printf("Existing users:")
	for _, client := range hub.clients {
		log.Println(client.RemoteAddr)
	}
}

// Broadcast a message to every client
func (hub *Hub) Broadcast(message []byte) {
	errors := stack.New()
	log.Printf("Broadcasting to %d clients", len(hub.clients))
	for index, client := range hub.clients {
		err := client.Send(message)
		if err != nil {
			log.Printf("Error sending to %s: %s\n", client.RemoteAddr, err)
			errors.Push(index)
		}
		log.Printf("Sent to %s\n", client.RemoteAddr)
	}

	// Prune old clients
	for errors.Len() > 0 {
		tempIndex := errors.Pop()
		if tempIndex == nil {
			break
		}
		index := tempIndex.(int)

		hub.DeleteClient(index)
	}
}

func (hub *Hub) DeleteClient(index int) {
	if len(hub.clients) == index {
		hub.clients = hub.clients[:index]
	} else if len(hub.clients) > index {
		hub.clients = append(hub.clients[:index], hub.clients[index+1:]...)
	}
}

var currentDebt = &messages.TechDebt{}

func checker(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: checker,
}

var hub = New()

func upgrade(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	//defer c.Close()
	log.Printf("Got connection from %s", r.RemoteAddr)

	// Immediately send current tech debt
	var buffer []byte
	buffer, err = proto.Marshal(currentDebt)
	if err != nil {
		log.Println("marshal err:", err)
		return
	}
	err = c.WriteMessage(websocket.BinaryMessage, buffer)
	if err != nil {
		log.Println("write err:", err)
		return
	}

	hub.AddClient(r.RemoteAddr, c)
}

func main() {

	flag.Parse()
	log.SetFlags(0)
	http.Handle("/", http.FileServer(rice.MustFindBox("client/dist").HTTPBox()))
	http.HandleFunc("/ws", upgrade)
	go handleAllUsers(&hub)
	log.Println("Upgrade using ws://127.0.0.1:3000/ws")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handleAllUsers(h *Hub) {
	currentDebt.TechDebt = 0
	var receiveMsg = &messages.TechDebt{}
	for {
		// Get message
		input := <-h.Receiver

		// Unmarshal message
		err := proto.Unmarshal(input, receiveMsg)
		if err != nil {
			log.Println("unmarshal err:", err)
			break
		}

		// Add new tech debt
		currentDebt.TechDebt += receiveMsg.TechDebt
		currentDebt.Username = receiveMsg.Username
		currentDebt.Color = receiveMsg.Color
		fmt.Printf("[%s]%s added %d points:\tCurrent Score: %d\n", receiveMsg.Color, receiveMsg.Username, receiveMsg.TechDebt, currentDebt.TechDebt)

		// Write tech debt
		buffer, err := proto.Marshal(currentDebt)
		if err != nil {
			log.Println("marshal err:", err)
			break
		}
		hub.Broadcast(buffer)
	}
}

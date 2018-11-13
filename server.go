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

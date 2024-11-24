package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string, 1)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var results string

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
}

// PostHandler converts post request body to string
func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		results = string(body)
		defer r.Body.Close()

		go func(res string) {
			broadcast <- res
		}(results)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// register client
	clients[ws] = true
}

func echo() {
	for {
		val := <-broadcast
		fmt.Println("FROM CHANNEL: " + val)
		// send to every client that is currently connected
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(val))
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/messages", PostHandler)
	http.HandleFunc("/ws", WsHandler)

	go echo()

	log.Println("Starting server on port " + os.Getenv("WORKER_PORT"))
	http.ListenAndServe(":"+os.Getenv("WORKER_PORT"), nil)
}

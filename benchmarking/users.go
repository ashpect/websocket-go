package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	serverAddr = "ws://localhost:8080" // WebSocket server address
	numClients = 50000                 // Number of concurrent WebSocket clients
	message    = "Hello"               // Message to send
)

func benchmarkClient(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Printf("Client %d: Failed to connect to server: %v", id, err)
		return
	}
	defer conn.Close()

	// Send the message
	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Client %d: Failed to send message: %v", id, err)
		return
	}

	log.Printf("Client %d successfully sent message: %s", id, message)
}

func main() {
	var wg sync.WaitGroup

	// Launch multiple concurrent clients
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go benchmarkClient(i, &wg)
	}

	// Wait for all clients to finish
	wg.Wait()
	fmt.Println("Benchmark finished")
}

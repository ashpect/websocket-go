package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	serverAddr = "ws://localhost:8080" // WebSocket server address
	numClients = 1000                  // Number of concurrent WebSocket clients
	message    = "Hello"               // Message to send
)

var (
	failedConnections     int
	failedSends           int
	successfulConnections int
	successfulSends       int
	totalBytesSent        int
)

func benchmarkClient(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Connect to the WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Printf("Client %d: Failed to connect to server: %v", id, err)
		// Increment failed connection counter
		failedConnections++
		return
	}
	defer conn.Close()

	// Increment successful connection counter
	successfulConnections++

	// Send the "Hello" message
	messageBytes := []byte(message)
	err = conn.WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		log.Printf("Client %d: Failed to send message: %v", id, err)
		// Increment failed send counter
		failedSends++
		return
	}

	// Increment successful send counter
	successfulSends++
	// Update total bytes sent
	totalBytesSent += len(messageBytes)
	log.Printf("Client %d successfully sent message: %s", id, message)
}

func main() {
	var wg sync.WaitGroup

	// Record start time for benchmarking
	startTime := time.Now()

	// Launch multiple concurrent clients
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go benchmarkClient(i, &wg)
	}

	// Wait for all clients to finish
	wg.Wait()

	// Calculate the duration
	duration := time.Since(startTime)
	fmt.Printf("Benchmark finished: %v\n", duration)

	// Calculate Throughput (Messages per second)
	messageThroughput := float64(successfulSends) / duration.Seconds()

	// Calculate Data Throughput (Bytes per second)
	dataThroughput := float64(totalBytesSent) / duration.Seconds()

	// Print the results
	totalConnections := numClients
	totalSends := numClients
	fmt.Printf("Total Connections: %d\n", totalConnections)
	fmt.Printf("Successful Connections: %d\n", successfulConnections)
	fmt.Printf("Failed Connections: %d\n", failedConnections)
	fmt.Printf("Total Messages Sent: %d\n", totalSends)
	fmt.Printf("Successful Messages Sent: %d\n", successfulSends)
	fmt.Printf("Failed Messages Sent: %d\n", failedSends)
	fmt.Printf("Message Throughput: %.2f messages/second\n", messageThroughput)
	fmt.Printf("Data Throughput: %.2f bytes/second\n", dataThroughput)
}

package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	serverAddr    = "ws://localhost:8080"
	numClients    = 50000
	message       = "Hello"
	retryMax      = 3
	retryWaitTime = 3 * time.Second
	sleepDuration = 10 * time.Millisecond // Delay to test various throughput
)

var (
	failedConnections     int
	failedSends           int
	successfulConnections int
	successfulSends       int
	totalBytesSent        int
	mu                    sync.Mutex
)

func benchmarkClient(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	var conn *websocket.Conn
	var err error

	for i := 0; i < retryMax; i++ {
		conn, _, err = websocket.DefaultDialer.Dial(serverAddr, nil)
		if err != nil {
			log.Printf("Client %d: Failed to connect to server: %v", id, err)

			time.Sleep(3 * time.Second)
			continue
		}
		break
	}

	if conn == nil {
		mu.Lock()
		failedConnections++
		mu.Unlock()
		return
	}

	defer conn.Close()
	mu.Lock()
	successfulConnections++
	mu.Unlock()

	// Send message
	messageBytes := []byte(message)
	err = conn.WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		log.Printf("Client %d: Failed to send message: %v", id, err)
		mu.Lock()
		failedSends++
		mu.Unlock()
		return
	}

	mu.Lock()
	successfulSends++
	mu.Unlock()

	_, inMsg, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Client %d: Failed to read message: %v", id, err)
	}
	totalBytesSent += len(messageBytes)
	log.Printf(string(inMsg))
}

func main() {
	var wg sync.WaitGroup

	startTime := time.Now()

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		// Staggering connections to make real world scenario
		time.Sleep(sleepDuration)
		go benchmarkClient(i, &wg)
	}

	wg.Wait()

	duration := time.Since(startTime)
	fmt.Printf("Benchmark finished: %v\n", duration)

	messageThroughput := float64(successfulSends) / duration.Seconds()
	dataThroughput := float64(totalBytesSent) / duration.Seconds()

	fmt.Printf("Total Connections: %d\n", numClients)
	fmt.Printf("Successful Connections: %d\n", successfulConnections)
	fmt.Printf("Failed Connections: %d\n", failedConnections)
	fmt.Printf("Total Messages Sent: %d\n", numClients)
	fmt.Printf("Successful Messages Sent: %d\n", successfulSends)
	fmt.Printf("Failed Messages Sent: %d\n", failedSends)
	fmt.Printf("Message Throughput: %.2f messages/second\n", messageThroughput)
	fmt.Printf("Data Throughput: %.2f bytes/second\n", dataThroughput)
}

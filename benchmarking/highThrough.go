package main

import (
	"fmt"
	"log"
	"runtime"
	"sort"
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
	sleepDuration = 5 * time.Millisecond // Delay to test various throughput
)

var (
	failedConnections     int
	failedSends           int
	successfulConnections int
	successfulSends       int
	totalBytesSent        int
	mu                    sync.Mutex
	latencies             []time.Duration
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

	// Connection started
	start := time.Now()
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

	// Connection ended, msg recieved.
	latency := time.Since(start)
	mu.Lock()
	latencies = append(latencies, latency)
	mu.Unlock()

	totalBytesSent += len(messageBytes)
	log.Printf(string(inMsg))
}

func memoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

func main() {
	var wg sync.WaitGroup

	startTime := time.Now()

	go func() {
		for {
			fmt.Printf("Memory Usage: %d KB\n", memoryUsage()/1024)
			time.Sleep(1 * time.Second)
		}
	}()

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

	// latency statistics
	var totalLatency time.Duration
	for _, lat := range latencies {
		totalLatency += lat
	}
	averageLatency := totalLatency / time.Duration(len(latencies))
	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	minLatency := latencies[0]
	maxLatency := latencies[len(latencies)-1]
	medianLatency := latencies[len(latencies)/2]

	fmt.Printf("Total Connections: %d\n", numClients)
	fmt.Printf("Successful Connections: %d\n", successfulConnections)
	fmt.Printf("Failed Connections: %d\n", failedConnections)
	fmt.Printf("Total Messages Sent: %d\n", numClients)
	fmt.Printf("Successful Messages Sent: %d\n", successfulSends)
	fmt.Printf("Failed Messages Sent: %d\n", failedSends)
	fmt.Printf("Message Throughput: %.2f messages/second\n", messageThroughput)
	fmt.Printf("Data Throughput: %.2f bytes/second\n", dataThroughput)

	// Latency statistics
	fmt.Printf("Latency (ms) - Avg: %.2f, Median: %.2f, Min: %.2f, Max: %.2f\n",
		float64(averageLatency.Milliseconds()),
		float64(medianLatency.Milliseconds()),
		float64(minLatency.Milliseconds()),
		float64(maxLatency.Milliseconds()),
	)

	finalMemoryUsage := memoryUsage()
	fmt.Printf("Final Memory Usage: %d KB\n", finalMemoryUsage/1024)
}

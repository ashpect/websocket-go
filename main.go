package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

type client struct {
	conn          *websocket.Conn
	writeMessage  chan string
	serverMessage chan string
	sessionID     uuid.UUID
	counter       int
}

type controller struct {
	clients  map[uuid.UUID]*client
	register chan *client
	delete   chan *client
	mu       sync.RWMutex
}

func (m *controller) run() {
	// fmt.Println("Running controller")
	for {
		select {
		case c := <-m.register:
			// fmt.Println("Registering client")
			m.mu.Lock()
			// fmt.Println("a")
			m.clients[c.sessionID] = c
			m.mu.Unlock()
		case c := <-m.delete:
			m.mu.Lock()
			delete(m.clients, c.sessionID)
			m.mu.Unlock()
		}
	}
}

func (m *controller) getClients() map[uuid.UUID]*client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.clients
}

func (m *controller) getClient(sessionID uuid.UUID) (*client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.clients[sessionID]
	if !ok {
		log.Printf("Client with session id %s not found", sessionID)
		return nil, fmt.Errorf("Client with session id %s not found", sessionID)
	}
	return a, nil
}

func (wsh *webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, controller *controller) {

	authHeader := r.Header.Get("Authorization")

	c := &client{}

	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		fmt.Println("Creating new client")
		c.conn = conn
		c.writeMessage = make(chan string)
		c.serverMessage = make(chan string)
		c.sessionID = uuid.New()
		fmt.Println("debuf")
		controller.register <- c
		fmt.Println("de")

	} else {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		parsedToken, err := validateToken(tokenString)
		if err != nil {
			log.Printf("Error %s when validating JWT token", err)
			return
		}
		claims := parsedToken.Claims.(jwt.MapClaims)
		parsedUUID, err := uuid.Parse(claims["user_id"].(string))
		fmt.Println("Parsed UUID is ", parsedUUID)
		if err != nil {
			log.Fatalf("Failed to parse user_id as UUID: %v", err)
			return
		}
		clientFound, err := controller.getClient(parsedUUID)
		if err != nil {
			response := "Client with session id " + parsedUUID.String() + " not found"
			err := conn.WriteMessage(websocket.TextMessage, []byte(response))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
				return
			}
			err = conn.Close()
			if err != nil {
				log.Printf("Error %s when closing connection", err)
			}
			return
		}
		fmt.Println("Client with session id " + parsedUUID.String() + "found")
		c = clientFound
		c.conn = conn
	}

	fmt.Println("Client session id is ", c.sessionID)
	go c.clientRead()
	// I dont need to handle stopping clientRead goroutine - it will automatically exit after the error reading msg when i stop clientHandler.
	go c.clientHandler()
}

func (c *client) clientHandler() {
	defer c.conn.Close()

	fmt.Println("Client handler started")
	token, err := createJWTToken(c.sessionID.String())
	if err != nil {
		log.Printf("Error %s when creating JWT token", err)
		return
	}

	reponse := "Connection Successful with session id " + c.sessionID.String() + ". Welcome to the server!. Your JWT token for futhur login is " + token
	c.clientWrite(reponse)

	// initializeTime := time.Now().Unix()
	// for (time.Now().Unix() - initializeTime) < 300 {
	sessionTimer := time.NewTimer(5 * time.Minute)
	defer sessionTimer.Stop() // Stop the timer when the function exits

	// enable server side pushing
	// Reads msgs
	for {
		select {
		case <-sessionTimer.C:
			response := "Session expired. Closing connection with client"
			c.clientWrite(response)
			return
		case serverMessage := <-c.serverMessage:
			c.clientWrite(serverMessage)
		case writeMessage := <-c.writeMessage:
			if string(writeMessage) == "close" {
				response := "Closing connection with client"
				c.clientWrite(response)
				return
			}
			c.counter++
			response := "Received message " + string(writeMessage) + " from client. This is message number " + fmt.Sprintf("%d", c.counter)

			c.clientWrite(response)
		}
	}

}

func (c *client) clientRead() {
	fmt.Println("Client read started")
	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error %s when reading message from client", err)
			return
		}
		if messageType != websocket.TextMessage {
			response := "Only text messages are supported"
			c.clientWrite(response)
			return
		}
		c.writeMessage <- string(message)
	}
}

func (c *client) clientWrite(response string) {
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(response))
	if err != nil {
		log.Printf("Error %s when sending message to client", err)
		return
	}
}

func main() {
	// Create a new controller
	controller := controller{
		clients:  make(map[uuid.UUID]*client),
		register: make(chan *client),
		delete:   make(chan *client),
		mu:       sync.RWMutex{},
	}
	go controller.run()

	webSocketHandler := webSocketHandler{
		upgrader: websocket.Upgrader{},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		webSocketHandler.ServeHTTP(w, r, &controller)
	})

	http.HandleFunc("/getClients", func(w http.ResponseWriter, r *http.Request) {
		clients := controller.getClients()
		for _, client := range clients {
			fmt.Println(client.sessionID)
		}
	})

	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

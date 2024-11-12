package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

type client struct {
	conn      *websocket.Conn
	serverMessage chan string
	sessionID uuid.UUID
	counter   int
}

type controller struct {
	clients  map[uuid.UUID]*client
	register chan *client
	delete   chan *client
	mu       sync.Mutex
}

func (m *controller) run() {
	for {
		select {
		case c := <-m.register:
			m.mu.Lock()
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
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.clients
}

func (m *controller) getClient(sessionID uuid.UUID) *client {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.clients[sessionID]
}

func (wsh *webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, controller *controller) {

	authHeader := r.Header.Get("Authorization")

	c := &client{}

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		conn, err := wsh.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("error %s when upgrading connection to websocket", err)
			return
		}
		c.conn = conn
		c.sessionID = uuid.New()
		fmt.Println("Creating new client")
		controller.register <- c
	} else {
		// get the old client session
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		parsedToken, err := validateToken(tokenString)
		if err != nil {
			log.Printf("Error %s when validating JWT token", err)
			return
		}
		fmt.Println(parsedToken.Claims)
		// TODO : extract the session id and then get the client from the controller
	}

	go c.clientHandler()

}

func (c *client) clientHandler() {
	conn := c.conn
	defer conn.Close()

	token, err := createJWTToken(c.sessionID.String())
	if err != nil {
		log.Printf("Error %s when creating JWT token", err)
		return
	}

	reponse := "Connection Successful with session id " + c.sessionID.String() + ". Welcome to the server!. Your JWT token for futhur login is " + token
	err = conn.WriteMessage(websocket.TextMessage, []byte(reponse))
	if err != nil {
		log.Printf("Error %s when sending message to client", err)
		return
	}

	// initializeTime := time.Now().Unix()
	// for (time.Now().Unix() - initializeTime) < 300 {
	for {

		// WebSocket protocol supports both text and binary message types. Since poc, only implementing for string types.
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error %s when reading message from client", err)
			return
		}

		if messageType != websocket.TextMessage {
			err = conn.WriteMessage(websocket.TextMessage, []byte("Only text messages are supported"))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
				return
			}
			return
		}

		if string(message) == "close" {
			err = conn.WriteMessage(websocket.TextMessage, []byte("Closing connection"))
			if err != nil {
				log.Printf("Error %s when sending message to client", err)
				return
			}
			return
		}

		c.counter++
		response := "Received message " + string(message) + " from client. This is message number " + fmt.Sprintf("%d", c.counter)
		err = conn.WriteMessage(websocket.TextMessage, []byte(response))
		if err != nil {
			log.Printf("Error %s when sending message to client", err)
			return
		}

	}

}

func main() {
	// Create a new controller
	controller := controller{
		clients: make(map[uuid.UUID]*client),
	}
	go controller.run()

	webSocketHandler := webSocketHandler{
		upgrader: websocket.Upgrader{},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		webSocketHandler.ServeHTTP(w, r, &controller)
	})

	http.HandleFunc("/getClients", func(w http.ResponseWriter, r *http.Request) {
		controller.mu.Lock()
		defer controller.mu.Unlock()
		for _, c := range controller.clients {
			fmt.Fprintf(w, "Session ID: %s, Counter: %d\n", c.sessionID, c.counter)
		}
	})

	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

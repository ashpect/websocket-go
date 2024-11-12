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

func (wsh *webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
		fmt.Println("Creating new client Connection")
		// controller.register <- c
	} else {
		// get the old client session
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		parsedToken, err := validateToken(tokenString)
		if err != nil {
			log.Printf("Error %s when validating JWT token", err)
			return
		}
		fmt.Println(parsedToken.Claims)
	}

	go clientHandler(c)

}

func clientHandler(c *client) {
	conn := c.conn
	defer conn.Close()

	token, err := createJWTToken(c.sessionID.String())
	if err != nil {
		log.Printf("Error %s when creating JWT token", err)
		return
	}

	reponse := "Connection Successful with session id " + c.sessionID.String() + ". Welcome to the server!. Your JWT token is " + token
	err = conn.WriteMessage(websocket.TextMessage, []byte(reponse))
	if err != nil {
		log.Printf("Error %s when sending message to client", err)
		return
	}

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

	http.Handle("/", &webSocketHandler)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

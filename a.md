type Upgrader struct {
    HandshakeTimeout time.Duration
    ReadBufferSize, WriteBufferSize int
    Subprotocols []string
    Error func(w http.ResponseWriter, r *http.Request, status int, reason error)
    CheckOrigin func(r *http.Request) bool
    EnableCompression bool
}


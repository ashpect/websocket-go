## Functional Requirements
1. Connection Management: <br>
[x] The server should allow users to connect and disconnect.
[x] Upon connection, the server should issue a unique session ID to the user. 
2. Message Handling: <br>
[x] While connected, users can send messages to the server.
[x] For each message received, the server should return a hardcoded response
along with the current message count for that session.
3. Server-Side Events: <br>
[x] The server should be capable of sending messages to clients without receiving a
client message (server-side push).
4. Benchmarking: <br>
[] Provide a basic benchmarking script that simulates clients. Each simulated client
should send the same dummy message and receive a response.
[] The benchmarking script must measure:
    []Total failures (e.g., no response received or connection dropped).
    []Total throughput of the WebSocket server.

## Non-Functional Requirements
1. Scalability: <br>
[] The server should handle at least 50,000 concurrent connections.
2. Deployment: <br>
[] Depoyment scripts in bash
[] Contanerize in docker/k8s


## Good-to-Have Features
1. Session Reconnection: <br>
[x] After disconnection, clients can reconnect to their session using the same session ID.
[x] Authtoken based reconnection and not simply session based.
2. Session Inactivity Management: <br>
[x] If a session is inactive for more than 5 minutes, the server should automatically
disconnect the user.
3. Advanced Benchmarking: <br>
[] In addition to the basic benchmarking:
[] Track and measure total memory usage.
[] Record the avg/median/low/high latency per session once a client
disconnects
[] Provide plots or visualizations of these metrics.
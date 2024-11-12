## Functional Requirements

1. Connection Management: <br>
[x] The server should allow users to connect and disconnect.
[x] Upon connection, the server should issue a unique session ID to the user. - verify
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

## Non-Functional Requirements

1. Scalability:
[] The server should handle at least 50,000 concurrent connections.
2. Deployment:
[] Depoyment scripts in bash
[] Contanerize in docker/k8s
## Good-to-Have Features

1. Session Reconnection:
[x] After disconnection, clients can reconnect to their session using the same session ID.
2. Session Inactivity Management:
[x] If a session is inactive for more than 5 minutes, the server should automatically
disconnect the user.
3. Advanced Benchmarking:
[] In addition to the basic benchmarking:
[] Track and measure total memory usage.
[] Record the avg/median/low/high latency per session once a client
disconnects
[] Provide plots or visualizations of these metrics.
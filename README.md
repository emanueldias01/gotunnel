# GoTunnel

GoTunnel is a TCP tunneling service written in Go. It allows a client application to expose a local TCP service through a public port on a GoTunnel server.

## How It Works

1. The server listens for tunnel connections on port `3000`.
2. A client connects to the server and opens a Yamux multiplexed session.
3. The server allocates an available public port from `9000` through `9999`.
4. The server sends the allocated port to the client over a control stream.
5. When a user connects to the public port, the server opens a Yamux stream and forwards the connection to the client's local application.
6. Closing the tunnel releases the public port.

```text
Internet client
      |
      | TCP
      v
Server: public port (9000-9999)
      |
      | Yamux multiplexed session
      v
GoTunnel client
      |
      | TCP
      v
Local application (for example, localhost:8080)
```

## Requirements

- Go 1.26 or later
- Network access to the server port `3000`
- A local TCP application for the client to forward traffic to

## Installation

Clone the repository and download dependencies:

```bash
git clone https://github.com/emanueldias01/gotunnel.git
cd gotunnel
go mod download
```

## Run the Server

Start the tunnel server:

```bash
go run ./cmd/server
```

The server listens on `0.0.0.0:3000` for client connections. It also binds each allocated public port on all interfaces.

For each connected client, the server prints a message similar to:

```text
Client connected: <client-id> (port 9123)
```

The public port range is fixed at `9000` to `9999`.

## Run the Client

The client requires the tunnel server host, the tunnel server port, and the local application port:

```bash
go run ./cmd/client <server-host> <server-port> <local-port>
```

For example, to expose an application running on `localhost:8080` through a server located at `tunnel.example.com`:

```bash
go run ./cmd/client tunnel.example.com 3000 8080
```

The client prints the public port assigned by the server:

```text
Open port 9123
```

You can then access the local application through:

```text
tunnel.example.com:9123
```

## Local Example

From the repository root, start the included test API:

```bash
go run ./cmd/test
```

The test API listens on `localhost:8080` and provides:

```text
GET /ping
```

Start the tunnel server in another terminal:

```bash
go run ./cmd/server
```

Start the client in a third terminal:

```bash
go run ./cmd/client localhost 3000 8080
```

If the server assigns port `9000`, test the tunnel with:

```bash
curl http://localhost:9000/ping
```

The response should be:

```json
{"message":"pong"}
```

## Project Structure

```text
cmd/
  client/    Tunnel client command
  server/    Tunnel server command
  test/      Local HTTP API used for testing
internal/
  client/    Server connection and local forwarding
  common/     Bidirectional TCP proxy
  server/     Client management and public port allocation
```

## Implementation Notes

- The Go Tunnel protocol is built on [Yamux](https://github.com/hashicorp/yamux).
- The server creates one multiplexed session per client and one stream per visitor connection.
- Public ports are allocated sequentially within the configured range and released when the client disconnects.
- The public ports and tunnel server port are bound to all interfaces.
- The client forwards traffic to `localhost` on the supplied local port.
- The project does not currently include authentication, encryption, TLS, or access control.

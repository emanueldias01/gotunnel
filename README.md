# GoTunnel

GoTunnel exposes a local TCP service through a port on a remote server. It uses [Yamux](https://github.com/hashicorp/yamux) to multiplex tunnel traffic.

## Install

Requires Go 1.26.4 or later.

```bash
go install github.com/emanueldias01/gotunnel/cmd/gotunnel@latest
```

Ensure `$(go env GOBIN)` (or `$(go env GOPATH)/bin` when `GOBIN` is unset) is in your `PATH`.

## Use

Start the tunnel server:

```bash
gotunnel server
```

The server listens on port `3000` and assigns each connected client an available port from `9000` to `9999`. Ensure these ports are reachable through the server's firewall. The server prints its local IPv4 addresses at startup.

On the machine running the application, connect the tunnel client:

```bash
gotunnel client <server-host:3000> <local-port>
```

For example, to expose an application listening on port `8080`:

```bash
gotunnel client tunnel.example.com:3000 8080
```

The client prints the server IP and assigned public port. Connect to `<server-host>:<assigned-port>` to reach the local application.

## Local test

In separate terminals, run the included test API, the tunnel server, and a client:

```bash
gotunnel test
gotunnel server
gotunnel client localhost:3000 8080
```

The test API listens on port `8080`. Use the port printed by the client to test it, for example:

```bash
curl http://localhost:9000/ping
```

`gotunnel test` is a sample HTTP service; it responds to `GET /ping` with `{"message":"pong"}`.

## Security

The current tunnel has no authentication or TLS. Do not expose it to untrusted networks without adding appropriate protections.

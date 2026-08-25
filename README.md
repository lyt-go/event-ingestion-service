# Event Ingestion Service

A small Go HTTP service for receiving events, managing subscriptions, delivering matching events, and publishing immutable stream snapshots. It uses an in-memory store so it can run without external services.

## Run

```bash
go run ./cmd/server
```

The server listens on `:8080` by default. Set `ADDR` to change the address.

## API

- `POST /v1/events` receives an event with a stream, type, and payload.
- `POST /v1/subscriptions` registers a subscriber for a stream and event type.
- `GET /v1/deliveries?subscriber=...` reads pending deliveries.
- `POST /v1/deliveries/{id}/ack` acknowledges a delivery.
- `POST /v1/snapshots` publishes a stream snapshot.
- `GET /v1/snapshots/{stream}` reads the latest snapshot.

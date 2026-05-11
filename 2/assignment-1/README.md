# Assignment 4 - Caching Strategies & Background Jobs

> small 3-service platform composed of a **Doctor Service**, an **Appointment Service**, and a **Notification Service**

## License

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/license/mit/)

## Project Overview

### Technology Stack
- Language: Go 1.26
- Framework: gRPC
- Database: PostgreSQL
- Configuration: TOML / Environment Variables

### What Changed from Assignment 3

- Added Redis cache layer for read endpoints.
- Added cache invalidation on write endpoints.
- Added Redis-backed per-client rate limiting via gRPC unary interceptors.
- Extended Notification Service with a worker-pool background queue.
- Added retry with exponential backoff for gateway failures.
- Added idempotency key storage in Redis to prevent duplicate external calls.
- Added new `mock-gateway` binary to simulate third-party notification API behavior.

## Caching Strategy
Redis cache is implemented behind cache repository interfaces and injected into services.
No Redis calls are made in domain models or gRPC handlers.

### Doctor Service Cache

| Operation | Strategy | Key | TTL |
|---|---|---|---|
| `GetDoctor(id)` | Cache-Aside | `doctor:<id>` | `CACHE_TTL_SECONDS` |
| `ListDoctors` | Cache-Aside | `doctors:list` | `CACHE_TTL_SECONDS` |
| `CreateDoctor` | Write-Through | `doctors:list` | immediate eviction |

### Appointment Service Cache

| Operation | Strategy | Key | TTL |
|---|---|---|---|
| `GetAppointment(id)` | Cache-Aside | `appointment:<id>` | `CACHE_TTL_SECONDS` |
| `ListAppointments` | Cache-Aside | `appointments:list` | `CACHE_TTL_SECONDS` |
| `CreateAppointment` | Write-Around | `appointments:list` | immediate eviction |
| `UpdateAppointmentStatus` | Write-Through | `appointment:<id>`, `appointments:list` | immediate eviction |

### Cache Invalidation Rules

- Invalidation occurs **after successful DB write** and **before gRPC response returns**.
- Cache miss never fails the request; DB is fallback source of truth.
- Cache write/delete failures are logged but do not block responses.
- Redis unavailability should degrade behavior to DB-only reads (best effort).

## Rate Limiting Algorithm

Rate limiter is implemented as `UnaryServerInterceptor` in Doctor and Appointment services.

Sliding window per minute counter using Redis:
- Key format: `<service>:ratelimit:<client_ip>:<minute_window>`
- For each request:
	1. `INCR` key
	2. set `EXPIRE 60s` when counter is first created
- If counter > configured RPM -> return `codes.ResourceExhausted`

Why this algorithm:
- Simple and deterministic for defense/demo.
- Centralized Redis counters keep limits consistent across multiple service instances.

## Background Job Queue

Notification service has 2 separated concerns:
- `internal/subscriber` - broker subscription
- `internal/jobqueue` - worker pool, idempotency, retry, dead-letter behavior

### Worker Pool
- Queue: buffered Go channel
- Configurable worker count (`WORKER_POOL_SIZE`, default 3)
- Hardcoded queue buffer
- Backpressure handling: enqueue blocks when channel is full (no silent drop)

### Idempotency
- Redis key: `notification:idempotency:<idempotency_key>`
- Value `done` means already processed
- TTL: 24h
- If already done -> job dropped (duplicate protection)

### Retry & Dead Letter
- Retries: up to 3 attempts
- Backoff: 1s, 2s, 4s
- Retry on:
  - HTTP 503 from gateway
  - network/unreachable gateway errors
- After 3 failures:
  - Write dead-letter JSON entry to `stderr`
  - Worker continues running

## How to Run the Project

```sh
docker compose up --build
```

## Startup Order

Using Docker Compose, dependencies and healthchecks are configured.

## Cache Consistency Trade-offs

- When Redis is unavailable, services should continue running and serve reads directly from PostgreSQL; caching is best-effort optimization, not a dependency for correctness.
- Redis is treated as performance optimization, not source of truth.
- If Redis fails, system still serves from PostgreSQL (higher latency).
- In distributed Redis/cluster mode, consistency and failover behavior depend on replication topology.

## Rate-Limiting Trade-offs

Limitation of local in-memory limiting in scaled deployments is per-instance counters are inconsistent across replicas -> clients can bypass limits by hitting different instances.

Redis-backed centralized counters solve both by sharing one global counter state.
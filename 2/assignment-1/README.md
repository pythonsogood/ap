# Assignment 3 - Message Queue & Database Migrations

> small 3-service platform composed of a **Doctor Service**, an **Appointment Service**, and a **Notification Service**

## License

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/license/mit/)

## Project Overview

### Technology Stack
- Language: Go 1.26
- Framework: gRPC
- Database: PostgreSQL
- Configuration: TOML / Environment Variables

### What Changed from Assignment 2

- Database: ~~SQLite~~ -> PostgreSQL
- SQL migrations (`golang-migrate`) for both services
- Event publishing after successful writes
- Notification service that subscribes to events and logs structured sJSON

### Broker Choice

Chosen broker: **NATS Core**.

Why: simple local setup and low overhead.

## Service Responsibilities

### Doctor Service
- Create new doctor profiles
- Retrieve doctor by ID
- List all doctors
- Enforce unique email constraint

### Appointment Service
- Create new appointments
- Retrieve appointment by ID
- List all appointments
- Update appointment status
- Validate doctor existence via Doctor Service gRPC

### Notification Service

- Connect to broker on startup.
- Retry broker connection with exponential backoff if unavailable.
- Subscribe to:
  - `doctors.created`
  - `appointments.created`
  - `appointments.status_updated`
- On each message, deserialize JSON and log one structured JSON message

## Environment Variables

| Service | Environment variables |
|---------|-----------------------|
| Shared | `NATS_CONNECTION_URL` -> `MESSAGE_BROKER_NATS_CONNECTION_URL` |
| Doctor | `DOCTOR_SERVICE_POSTGRES_USER` -> `DB_POSTGRES_USER`, `DOCTOR_SERVICE_POSTGRES_PASSWORD` -> `DB_POSTGRES_PASSWORD`, `DOCTOR_SERVICE_POSTGRES_DB` -> `DB_POSTGRES_DB`, `DB_POSTGRES_HOST`, `DB_POSTGRES_PORT` |
| Appointment | `APPOINTMENT_SERVICE_POSTGRES_USER` -> `DB_POSTGRES_USER`, `APPOINTMENT_SERVICE_POSTGRES_PASSWORD` -> `DB_POSTGRES_PASSWORD`, `APPOINTMENT_SERVICE_POSTGRES_DB` -> `DB_POSTGRES_DB`, `DB_POSTGRES_HOST`, `DB_POSTGRES_PORT`, `DOCTOR_SERVICE_ADDRESS` -> `SERVICES_DOCTOR_ADDRESS`, `DOCTOR_SERVICE_TIMEOUT` -> `SERVICES_DOCTOR_TIMEOUT` |
| Notification | `NOTIFICATION_SERVICE_LOG_DIRECTORY` `NOTIFICATION_SERVICE_LOG_FILE` -> `LOG_FILE` |

## How to Run the Project

```sh
docker compose up --build
```

## Migrations

Migrations run automatically on service startup before gRPC server starts.

Migration files:
- `doctor-service/migrations/`
  - `000001_create_doctors.up.sql`
  - `000001_create_doctors.down.sql`
- `appointment-service/migrations/`
  - `000001_create_appointments.up.sql`
  - `000001_create_appointments.down.sql`

## Startup Order

Using Docker Compose, dependencies and healthchecks are configured.

NATS broker and PostgreSQL instances launch first
Notification service depends on NATS broker, starts after the broker
Doctor and Appointment services depends on both NATS broker and their PostgreSQL databases, starts after the broker and a own db

## Event Contract

Published events are JSON with at least `event_type`, `occurred_at`, and payload fields.

| Subject | Publisher | Trigger | Fields |
|---------|-----------|---------|--------|
| `doctors.created` | `doctor-service` | `CreateDoctor` | `event_type`, `occurred_at`, `id`, `full_name`, `specialization`, `email` |
| `appointments.created` | `appointment-service` | `CreateAppointment` | `event_type`, `occurred_at`, `id`, `title`, `doctor_id`, `status` |
| `appointments.status_updated` | `appointment-service` | `UpdateAppointmentStatus` | `event_type`, `occurred_at`, `id`, `old_status`, `new_status` |

Example payload:

```json
{
	"event_type": "appointments.created",
	"occurred_at": "2026-05-01T10:24:01Z",
	"id": "appt-1",
	"title": "Initial cardiac consultation",
	"doctor_id": "doc-1",
	"status": "new"
}
```

## Consistency Trade-offs

DB write can succeed while publish could fail.

How to improve reliability:
- Outbox pattern (store event in DB transaction + background publisher)
- NATS JetStream

## NATS Core vs RabbitMQ

| Broker | NATS Core | RabbitMQ |
|--------|-----------|----------|
| Delievery model | fast fire-and-forget pub/sub | exchange/queue |
| Durability | no persistence by default | durable queues + acknowledgements |
| When to choose | simple transient notifications | guaranteed delivery matter |
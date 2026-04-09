# Assignment 2 - gRPC Migration

> small 2-service platform composed of a **Doctor Service** and an **Appointment Service**

## License

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/license/mit/)

## Project Overview

### Technology Stack
- Language: Go 1.26
- Framework: gRPC
- Database: SQLite3
- Configuration: TOML

### Architecture Diagram
![architecture diagram](/2/assignment-1/assets/architecture-diagram-grpc.svg)

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

## Proto Contract Description

### Doctor Service (doctor.proto)

| RPC | Request | Response | Business Rule |
|-----|---------|----------|---------------|
| `GetDoctor` | `GetDoctorRequest { string id }` | `DoctorResponse { id, full_name, specialization, email }` | Returns doctor by ID, returns NotFound if not exists |
| `GetDoctors` | `GetDoctorsRequest {}` | `GetDoctorsResponse { repeated DoctorResponse doctors }` | Returns all doctors |
| `CreateDoctor` | `CreateDoctorRequest { full_name, specialization, email }` | `DoctorResponse` | Creates doctor, enforces unique email constraint |

### Appointment Service (appointment.proto)

| RPC | Request | Response | Business Rule |
|-----|---------|----------|---------------|
| `GetAppointment` | `GetAppointmentRequest { string id }` | `AppointmentResponse` | Returns appointment by ID, returns NotFound if not exists |
| `GetAppointments` | `GetAppointmentsRequest {}` | `GetAppointmentsResponse { repeated AppointmentResponse }` | Returns all appointments |
| `CreateAppointment` | `CreateAppointmentRequest { title, description, doctor_id }` | `AppointmentResponse` | Validates doctor_id via Doctor Service gRPC before creation |
| `UpdateAppointmentStatus` | `UpdateAppointmentStatusRequest { id, status }` | `UpdateAppointmentStatusResponse` | Updates status (NEW -> IN_PROGRESS -> DONE) |

**AppointmentStatus Enum:**
- `NEW` (0): Initial state when created
- `IN_PROGRESS` (1): Appointment is being conducted
- `DONE` (2): Appointment completed

## REST vs gRPC Trade-offs

| Aspect | REST | gRPC | When to Choose |
|--------|------|------|----------------|
| **Protocol** | HTTP/1.1 + JSON | HTTP/2 + Protocol Buffers | gRPC for performance, REST for browser compatibility |
| **Contract** | OpenAPI/Swagger (manual) | .proto files (code-generated) | gRPC for strong typing, REST for flexibility |
| **Code Generation** | Manual implementation | Auto-generated | gRPC for faster development, REST for more control |
| **Streaming** | requires WebSockets | Native | gRPC for real-time, REST for simpler use cases |
| **Browser Support** | Universal | Limited | REST for web clients, gRPC for internal services |

### Why gRPC for this project?
- Strong typing with Protocol Buffers reduces runtime errors
- Code generation ensures client/server contracts stay in sync
- HTTP/2 provides better performance
- Smaller payload sizes

## Folder Structure

```
┌───appointment-service
│   │   .dockerignore
│   │   Dockerfile
│   │   go.mod
│   │   go.sum
│   │
│   ├───cmd
│   │   └───appointment
│   │       │   appointment.go
│   │       │
│   │       └───config
│   │               config.go
│   │
│   ├───configs
│   │   └───appointment
│   │           config.toml
│   │
│   └───internal
│       ├───database
│       │       model.go
│       │       sqlite.go
│       │
│       ├───handler
│       │       appointment.go
│       │
│       ├───model
│       │       appointment.go
│       │
│       ├───repository
│       │       appointment.go
│       │
│       ├───service
│       │       appointment.go
│       │       doctor_service.go
│       │
│       └───transport
│           ├───grpc
│           │       appointment.go
│           │
│           └───http
│                   appointment.go
│
├───doctor-service
│   │   .dockerignore
│   │   Dockerfile
│   │   go.mod
│   │   go.sum
│   │
│   ├───cmd
│   │   └───doctor
│   │       │   doctor.go
│   │       │
│   │       └───config
│   │               config.go
│   │
│   ├───configs
│   │   └───doctor
│   │           config.toml
│   │
│   └───internal
│       ├───database
│       │       model.go
│       │       sqlite.go
│       │
│       ├───handler
│       │       doctor.go
│       │
│       ├───model
│       │       doctor.go
│       │
│       ├───repository
│       │       doctor.go
│       │
│       ├───service
│       │       doctor.go
│       │
│       └───transport
│           ├───grpc
│           │       doctor.go
│           │
│           └───http
│                   doctor.go
│
└───proto
    │   appointment.proto
    │   doctor.proto
    │
    └───go
        └───proto
                appointment.pb.go
                appointment_grpc.pb.go
                doctor.pb.go
                doctor_grpc.pb.go
                go.mod
                go.sum
```

## Dependency Flow
![dependency flow](/2/assignment-1/assets/dependency-flow-grpc.svg)

## Inter-Service Communication

### How Appointment Service Validates Doctor

The Appointment Service calls the Doctor Service via gRPC to validate that a doctor exists before creating an appointment.

**gRPC Contract:**
- Service: `DoctorService`
- Method: `GetDoctor(GetDoctorRequest) returns (DoctorResponse)`
- Timeout: `15 seconds`

Implementation ([appointment-service/internal/service/doctor_service.go](/2/assignment-1/appointment-service/internal/service/doctor_service.go)):
```go
func (s grpcDoctorServiceImpl) IsValidDoctorId(doctor_id string) (bool, error) {
	conn, err := grpc.NewClient(s.doctor_service_address, grpc.WithTransportCredentials(insecure.NewBundle().TransportCredentials()))

	if err != nil {
		return false, err
	}

	defer conn.Close()

	client := pb.NewDoctorServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	_, err = client.GetDoctor(ctx, &pb.GetDoctorRequest{
		Id: doctor_id,
	})

	if err != nil {
		if strings.Contains(err.Error(), "code = NotFound") {
			return false, nil
		}

		return false, fmt.Errorf("[ERROR] Doctors service is currently unavailable\nError: %s", err.Error())
	}

	return true, nil
}
```

**Error Handling:**
- If doctor not found (`NotFound`): returns `false, nil` -> appointment creation fails with validation error
- If service unavailable: returns error with `Unavailable` status -> appointment creation fails

## How to Run the Project

```sh
docker compose up --build
```

This will start both services:
- Doctor Service: localhost:8081
- Appointment Service: localhost:8082

## Why Separate Databases

Each service owns its data independently:
- Doctor Service owns the `doctors` table
- Appointment Service owns the `appointments` table

Pros of separate databases:
- Encapsulation: each service's data is only accessible through its gRPC
- Independence: services can develop independently

### Why Not a Shared Database?

A shared database would create a distributed monolith:
- Changes to one service's data model could break the other
- Tight coupling between services

## Failure Scenario

### When Doctor Service is Unavailable

If the Doctor Service is unreachable when creating/updating an appointment:

1. Appointment Service attempts gRPC call to Doctor Service
2. Context timeout triggers after `15 seconds` (configurable)
3. gRPC returns `Unavailable` status code
4. Error is returned to client with details

**gRPC Status Codes:**
- `Unavailable` (code 14): Doctors service is currently unavailable
- `NotFound` (code 5): Doctor with id <ID> not found

**Logged Output:**
```
[ERROR] Doctors service is currently unavailable
Error: <error_details>
```
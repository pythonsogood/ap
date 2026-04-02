# Assignment 1 - Clean Architecture-Based Microservices

> small 2-service platform composed of a **Doctor Service** and an **Appointment Service**

## License

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/license/mit/)

## Project Overview

### Technology Stack
- Language: Go 1.26
- Framework: Gin
- Database: SQLite3
- Configuration: TOML

### Architecture Diagram
![architecture diagram](/2/assignment-1/assets/architecture-diagram.svg)

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
- Validate doctor existence via Doctor Service API

## Folder Structure

```
appointment-service
│   .dockerignore
│   Dockerfile
│   go.mod
│   go.sum
│
├───cmd
│   └───appointment
│       │   appointment.go
│       │
│       └───config
│               config.go
│
├───configs
│   └───appointment
│           config.toml
│
└───internal
    ├───database
    │       model.go
    │       sqlite.go
    │
    ├───handler
    │       appointment.go
    │
    ├───model
    │       appointment.go
    │
    ├───repository
    │       appointment.go
    │
    ├───service
    │       appointment.go
    │       doctor_service.go
    │
    └───transport
        └───http
                appointment.go

doctor-service
│   .dockerignore
│   Dockerfile
│   go.mod
│   go.sum
│
├───cmd
│   └───doctor
│       │   doctor.go
│       │
│       └───config
│               config.go
│
├───configs
│   └───doctor
│           config.toml
│
└───internal
    ├───database
    │       model.go
    │       sqlite.go
    │
    ├───handler
    │       doctor.go
    │
    ├───model
    │       doctor.go
    │
    ├───repository
    │       doctor.go
    │
    ├───service
    │       doctor.go
    │
    └───transport
        └───http
                doctor.go
```

## Dependency Flow
![dependency flow](/2/assignment-1/assets/dependency-flow.svg)

## Inter-Service Communication

### How Appointment Service Validates Doctor

HTTP Contract:
- Endpoint: `GET http://doctor-service:8081/doctors/{doctor_id}`
- Timeout: `15 seconds`
- Success Response: HTTP 200 OK
- Failure Response: Any other status code

Implementation ([appointment-service/internal/service/doctor_service.go](/2/assignment-1/appointment-service/internal/service/doctor_service.go)):
```go
func (s doctorServiceImpl) IsValidDoctorId(doctor_id string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/doctors/%s", s.doctor_service_address, doctor_id), nil)

	if err != nil {
		return false, err
	}

	req.Header.Set("X-Service-Name", "appointment-service")

	resp, err := s.HTTPClient.Do(req)

	if err != nil {
		log.Println("[ERROR] Doctors service is currently unavailable\nError: ", err.Error())

		return false, fmt.Errorf("Doctors service is currently unavailable\nError: %s", err.Error())
	}

	return resp.StatusCode == http.StatusOK, nil
}
```

## How to Run the Project

```sh
docker compose up --build
```

This will start both services:
- Doctor Service: http://localhost:8081
- Appointment Service: http://localhost:8082

## Why Separate Databases

Each service owns its data independently:
- Doctor Service owns the `doctors` table
- Appointment Service owns the `appointments` table

Pros of separate databases:
- Encapsulation: each service's data is only accessible through its API
- Independence: services can develop independently

### Why Not a Shared Database?

A shared database would create a distributed monolith:
- Changes to one service's data model could break the other
- Tight coupling between services

## Failure Scenario

### When Doctor Service is Unavailable

If the Doctor Service is unreachable when creating/updating an appointment:
1. Appointment Service attempts HTTP call to Doctor Service
2. Timeout triggers after `15 seconds`
3. Error is returned to client
4. Failure is logged

Response to Client:
```json
{
	"error": "Doctors service is currently unavailable\nError: <error_details>"
}
```

Logged Output:
```
[ERROR] Doctors service is currently unavailable
Error: <error_details>
```

## API Examples

### Doctor Service (Port 8081)

#### Create Doctor
```sh
curl -X POST http://localhost:8081/doctors \
	-H "Content-Type: application/json" \
	-d '{
		"full_name": "Dr. Aisha Seitkali",
		"specialization": "Cardiology",
		"email": "a.seitkali@clinic.kz"
	}'
```

#### Get Doctor by ID
```sh
curl http://localhost:8081/doctors/{id}
```

#### List All Doctors
```sh
curl http://localhost:8081/doctors
```

### Appointment Service (Port 8082)

#### Create Appointment
```sh
curl -X POST http://localhost:8082/appointments \
	-H "Content-Type: application/json" \
	-d '{
		"title": "Initial cardiac consultation",
		"description": ""Patient referred for palpitations and shortness of breath",
		"doctor_id": "<doctor-id-from-above>"
	}'
```

#### Get Appointment by ID
```sh
curl http://localhost:8082/appointments/{id}
```

#### List All Appointments
```sh
curl http://localhost:8082/appointments
```

#### Update Appointment Status
```sh
curl -X PATCH http://localhost:8082/appointments/{id}/status \
	-H "Content-Type: application/json" \
	-d '{
		"status": "in_progress"
	}'
```
# Assignment 2: Generic Concurrent Web Server

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/license/mit/)

## Usage

run

```sh
go run .
```

the application will be available at: http://localhost:8000

## Endpoints

- `GET /tasks`, returns `[ {"id": string, "status": string} ]`
- `POST /tasks`, body `{"payload": string}`
- `GET /tasks/{id}`, returns `{"id": string, "status": string}`
- `GET /stats`, returns `{"submitted": int, "completed": int, "in_progress": int}`

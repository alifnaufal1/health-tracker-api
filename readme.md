# Health Tracker API

A personal REST API backend for a self-hosted health/activity tracker, built to replace a third-party smartwatch companion app with a system I fully own and control.

This backend is part of a two-repo personal project:

- **Backend (this repo):** [health-tracker-api](https://github.com/alifnaufal1/health-tracker-api) — Go / Fiber / PostgreSQL
- **Mobile app:** [health-tracker-app](https://github.com/alifnaufal1/health-tracker-app) — React Native, connects to the smartwatch over BLE and talks to this API

> **Status:** Active development (early stage / MVP). Core user & device management is implemented; workout/health-metric ingestion endpoints are planned next. See [Roadmap](#roadmap).

---

## Background

The stock companion app for my Haylou smartwatch is ad-heavy, has a confusing UX, and doesn't expose the raw sensor data in a useful way. This project is a self-directed learning exercise to:

- Own my activity/health data end-to-end on infrastructure I control.
- Get hands-on with BLE device integration on the mobile side.
- Practice building a clean, layered Go backend and containerized deployment.

## Features

- **User management** — create, update, delete, list, and fetch users by ID.
- **Device management** — register and manage smartwatch devices (name, serial number, firmware/software revision, manufacturer), linked to a user.
- **Consistent JSON response envelope** (`code`, `status`, `data`, `message`) across all endpoints.
- **Centralized error handling** with typed application errors (`NotFoundError`, `ValidationError`, `ConflictError`) mapped to proper HTTP status codes.
- **Request validation** via `go-playground/validator`.
- **Password hashing** with bcrypt.
- **Structured logging** with Logrus, plus per-request IDs via Fiber's `requestid` middleware.
- **JWT middleware** is implemented and ready (`pkg/middleware/auth.go`), to be wired onto protected routes as auth requirements solidify.
- **OpenAPI 3.0 spec** describing the intended public API surface ([`api-spec.yml`](./api-spec.yml)).

## Tech Stack

| Layer            | Technology                                   |
| ---------------- | -------------------------------------------- |
| Language         | Go 1.25                                      |
| Web framework    | [Fiber v3](https://github.com/gofiber/fiber) |
| ORM / Database   | GORM + PostgreSQL                            |
| Validation       | go-playground/validator                      |
| Auth (prepared)  | JWT (`gofiber/contrib/jwt`)                  |
| Logging          | Logrus                                       |
| Containerization | Docker & Docker Compose                      |

## Project Structure

```
health-tracker-api/
├── main.go                     # App entrypoint: wiring + server bootstrap
├── internal/
│   ├── config/                 # Environment variable access
│   ├── controller/              # HTTP handlers (User, Device)
│   ├── service/                 # Business logic
│   ├── repository/              # Data access layer (GORM)
│   ├── router/                  # Route registration
│   ├── model/
│   │   ├── domain/              # GORM entities (User, Device, Base)
│   │   └── web/                 # Request/response DTOs
│   └── helper/                  # Response helpers, hashing, logging, tx helpers
├── pkg/
│   ├── database/                # DB connection + auto-migration
│   ├── middleware/               # Global error handler, JWT auth middleware
│   ├── logger/                  # Logger setup
│   └── apperror/                # Typed application errors
├── api-spec.yml                 # OpenAPI 3.0 specification
├── Dockerfile
└── docker-compose.yml           # API + PostgreSQL, for local/dev/production use
```

The codebase follows a layered architecture (**controller → service → repository → domain**) to keep HTTP concerns, business logic, and data access separate and testable.

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose (recommended, no local Postgres install needed)

### Environment Variables

Create a `.env` file in the project root:

```env
DB_HOST=postgres-db
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=password
DB_NAME=health-tracker-db
SECRET=your-jwt-signing-secret
```

### Run with Docker Compose (recommended)

This spins up the API and a PostgreSQL container together:

```bash
docker compose up --build
```

The API will be available at `http://localhost:3108`.

### Run locally (without Docker)

Make sure PostgreSQL is running and reachable with the credentials in your `.env`, then:

```bash
go mod download
go run main.go
```

The server listens on port `3108` by default.

## API Overview

Full request/response schemas are defined in [`api-spec.yml`](./api-spec.yml) (OpenAPI 3.0) — open it in [Swagger Editor](https://editor.swagger.io/) to explore interactively.

| Method   | Endpoint          | Description                  |
| -------- | ----------------- | ---------------------------- |
| `POST`   | `/api/user`       | Create a user                |
| `GET`    | `/api/user`       | List all users               |
| `GET`    | `/api/user/:id`   | Get a user by ID             |
| `PUT`    | `/api/user/:id`   | Update a user                |
| `DELETE` | `/api/user/:id`   | Delete a user                |
| `POST`   | `/api/device`     | Register a device for a user |
| `GET`    | `/api/device`     | List all devices             |
| `GET`    | `/api/device/:id` | Get a device by ID           |
| `PUT`    | `/api/device/:id` | Update a device              |
| `DELETE` | `/api/device/:id` | Delete a device              |

All responses follow this shape:

```json
{
  "code": 200,
  "status": true,
  "data": {},
  "message": "..."
}
```

## Roadmap

- [ ] Workout/activity data model + ingestion endpoints (heart rate, steps, pace) sent from the mobile app after BLE decoding.
- [ ] Wire up JWT authentication on protected routes (middleware already exists).
- [ ] CI/CD pipeline for automated build and deployment.
- [ ] Deploy to a personal server as a stable, portfolio-ready deployment.

## Related

- Mobile app (React Native + BLE): [health-tracker-app](https://github.com/alifnaufal1/health-tracker-app)

## Author

**Alif Naufal Muhammad**
[GitHub](https://github.com/alifnaufal1) · [LinkedIn](https://www.linkedin.com/in/naufal-muhammad-174a1332a/)

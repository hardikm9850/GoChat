# GoChat

GoChat is a mobile‑first real‑time chat backend built with **Go**. The project is designed for learning purpose with **clean architecture**, **domain‑driven principles**, and **scalability** in mind.

---

## Key Features

* Real‑time messaging support
* JWT‑based authentication
* Clean architecture (domain‑first design)
* Versioned REST APIs
* WebSocket support for live messaging
* MySQL‑backed persistence
* Production deployment on Railway
* Mobile‑friendly API contracts (Android & iOS)

---
## API Documentation (Swagger)

GoChat exposes a fully documented REST API using **Swagger (OpenAPI)**.

### Live Swagger UI
The latest API documentation is available here:

🔗 https://gochat-production-0f79.up.railway.app/swagger/index.html

This UI provides:
- Complete endpoint listings
- Request/response schemas
- Authentication details
- Ability to try APIs directly from the browser

### Local Access
When running the server locally, Swagger UI is available at:
http://localhost:8080/swagger/index.html


> Note: The Swagger documentation is generated from code annotations and always reflects the current API implementation.

--- 
## Architecture Overview

GoChat follows **Clean Architecture**, ensuring that business logic is independent of frameworks, databases, and delivery mechanisms.

### Layers

* **Domain**
  Core business entities, value objects, and repository contracts. No external dependencies.

* **Use Cases (Application Layer)**
  Application‑specific business rules. Orchestrates domain logic and enforces workflows.

* **Infrastructure**
  Database implementations, external services, persistence models, and framework integrations.

* **Delivery (HTTP / WebSocket)**
  REST APIs, WebSocket handlers, middleware, and request/response mapping.

---

## TODO : High‑Level System Diagram


* API Gateway / HTTP layer
* Authentication flow (JWT)
* WebSocket real‑time messaging
* Use case orchestration
* Database interactions (MySQL)


---

## API Design


```
/api/v1/health
/api/v1/auth/*
/api/v1/conversations
/api/v1/messages
/api/v1/ws
```

### Health Check

```http
GET /api/v1/health
```

Returns service status and is used for deployment health checks.

---

## Authentication

* JWT‑based access tokens
* Refresh token support (TODO)
* Token validation middleware
* Designed for mobile‑first usage

---

## Deployment

The backend is deployed on **Railway** using GitHub Actions.

* Automatic deployment on push to `main`
* Environment‑based configuration
* Secure secret management via GitHub Secrets



---

## Client Applications (Planned)

* **Android**: Jetpack Compose, MVI, Offline‑first
* **iOS**: SwiftUI or Kotlin Multiplatform (TBD)

The backend API contracts are designed to support both platforms consistently.

---

## Tech Stack

* **Language**: Go
* **Framework**: Gin
* **Database**: MySQL
* **ORM**: GORM
* **Auth**: JWT
* **Real‑time**: WebSockets
* **Deployment**: Railway
* **swaggo/swag** Swagger Docs

---

## Project Structure (Simplified)

```
internal/
  auth/
  chat/
    domain/
    usecase/
    repository/
    handler/
  infrastructure/
cmd/
```

---

## Testing (Planned)

* Unit tests for domain and use cases
* Repository tests with test containers
* API contract tests

---

## License

MIT License

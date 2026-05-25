# 📦 Polyglot Commerce — Suppliers API

## 📌 Overview
The **Suppliers API** handles internal corporate fulfillment flows, concurrent inventory allocation safety, and supplier coordination mechanics across both standard retail pipelines and B2B wholesale orders.

* **Technology Stack:** Go (Golang) | PostgreSQL
* **Architectural Pattern:** Idiomatic Go Clean Architecture
* **Concurrency Model:** Uses native Go channels & mutex mechanisms for strict double-allocation race prevention.

---

## 📂 Project Directory Architecture
```text
suppliers-api/
├── cmd/
│   └── server/          # App setup wrapper and main execution target
├── config/              # Environment config loading structures
├── internal/
│   ├── delivery/http/   # HTTP REST Handlers & router bindings
│   ├── domain/          # Core Enterprise Models & Interface Contracts
│   ├── repository/      # Persistent SQL storage operations
│   └── usecase/         # Pure business logic implementation
├── go.mod               # Go module dependency management
└── README.md
🚀 Local Development Setup
1. Download Dependencies
Bash
go mod tidy
2. Required Environment Configurations (.env)
Code snippet
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=suppliers_db
SERVER_PORT=:8083
3. Build and Run Server
Bash
go run cmd/server/main.go

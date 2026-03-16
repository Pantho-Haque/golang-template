# Prism - Parcel Backend Service

Prism is a high-performance Go-based backend service that powers the Parcel delivery platform's multiple frontend dashboards.

## 🏗️ Project Structure

```
prism/
├── cmd/                          # Application entry points
│   └── api/
│       └── main.go              # API server entry point
├── internal/                     # Private application code
│   ├── api/                     # HTTP API layer
│   │   ├── handlers/           # Handlers for the API
│   │   │
│   │   ├── middleware/
│   │   │   └── hubble_middleware.go  # Hubble APM middleware
│   │   └── routes.go            # Route definitions and setup
│   ├── cnst/                    # Constants and enums
│   │
│   ├── config/                  # Configuration management
│   │
│   ├── conn/                    # External service and dependency connections
│   │
│   ├── models/                  # Data models and entities
│   │
│   ├── monitoring/              # Monitoring and observability
│   │
│   ├── providers/               # External http service providers
│   │
│   ├── services/                # Business logic services
│   │
│   └── stores/                  # Data access layer
│
├── pkg/                         # Public packages
│
├── vendor/                      # Go module dependencies (4376 files)
├── .gitignore                   # Git ignore rules
├── Dockerfile                   # Docker container configuration
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── run.sh                       # Development run script
└── README.md                    # This file
```

run-management-service/
├── api/
│   └── openapi/
│       └── spec.yaml        # Our OpenAPI specification
├── cmd/
│   └── server/
│       └── main.go         # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go       # Configuration structures and loading
│   ├── domain/
│   │   ├── run.go         # Core domain models
│   │   └── errors.go      # Domain-specific errors
│   ├── ports/
│   │   ├── http/
│   │   │   ├── handlers/  # HTTP handlers
│   │   │   ├── middleware/ # HTTP middleware
│   │   │   └── server.go  # HTTP server setup
│   │   └── repository/
│   │       ├── mongodb/   # MongoDB implementation
│   │       └── interfaces.go # Repository interfaces
│   └── service/
│       └── run/           # Business logic
│           ├── service.go
│           └── service_test.go
├── pkg/
│   └── generated/         # Generated code (from OpenAPI)
├── scripts/
│   └── generate.sh       # Code generation scripts
├── .gitignore
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
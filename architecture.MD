backend/
├── cmd/
│   └── http/
│       ├── main.go
│       └── routes/
│           └── routes.go
├── internal/
│   ├── handlers/
│   │    ├── user_handler.go
│   │    └── home_handler.go
│   ├── middlewares/
│   │    └── auth_middleware.go
│   ├── models/
│   │    └── user.go
│   ├── repositories/
│   │    └── user_repository.go
│   ├── utils/
│   │    └── hash.go
│   └── validacoes/
│        └── user_validation.go
└── go.mod

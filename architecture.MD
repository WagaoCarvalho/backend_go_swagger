.
├── Insomnia_2025-03-11.json
├── Makefile
├── architecture.MD
├── auth
│   └── auth.go
├── cmd
│   └── http
│       ├── main.go
│       └── routes       
│           └── routes.go
├── config
│   ├── config.go
│   ├── db.go
│   ├── jwt.go
│   └── server.go
├── db
│   └── migrations
│       ├── 000001_create_users_table.down.sql
│       ├── 000001_create_users_table.up.sql
│       ├── 000002_create_wallets_table.down.sql
│       ├── 000002_create_wallets_table.up.sql
│       ├── 000003_create_transactions_table.down.sql
│       └── 000003_create_transactions_table.up.sql
├── docker-compose.yaml
├── go.mod
├── go.sum
├── important.MD
└── internal
    ├── handlers
    │   ├── home_handler.go
    │   ├── login_handler.go
    │   ├── transaction_handler.go
    │   ├── user_handler.go
    │   └── wallets_handler.go
    ├── middlewares
    │   ├── cors.go
    │   ├── jwt.go
    │   ├── logging.go
    │   ├── rate_limiter.go
    │   └── recovery_panic.go
    ├── models
    │   ├── transaction.go
    │   ├── user.go
    │   └── wallet.go
    ├── repositories
    │   ├── db_conn.go
    │   ├── transaction.go
    │   ├── user.go
    │   └── wallet.go
    ├── utils
    │   ├── generate_public_key.go
    │   ├── hash.go
    │   └── json_utils.go
    └── validations
        ├── user_validation.go
        └── validation.go
backend/
├── cmd/
│   └── http/
│       ├── main.go               # Ponto de entrada da aplicação
│       └── routes/
│           └── routes.go         # Configuração das rotas
├── internal/                     # Pacotes internos (não exportáveis)
│   ├── handlers/                 # Lógica dos handlers (controladores)
│   │    ├── user_handler.go
│   │    └── home_handler.go
│   ├── middlewares/              # Middlewares da aplicação
│   │    └── auth_middleware.go
│   ├── models/                   # Definição de modelos/entidades
│   │    └── user.go
│   ├── repositories/             # Comunicação com o banco de dados
│   │    └── user_repository.go
│   └── utils/                    # Utilitários e funções auxiliares
│        └── hash.go
└── go.mod



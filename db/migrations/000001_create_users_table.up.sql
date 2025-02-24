CREATE TABLE IF NOT EXISTS users (
    uid BIGSERIAL PRIMARY KEY,
    nickname VARCHAR(15) NOT NULL UNIQUE,
    email VARCHAR(40) NOT NULL UNIQUE,
    passwd VARCHAR(255) NOT NULL, -- Considerando armazenamento de hash de senha
    status CHAR(1) DEFAULT '0',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
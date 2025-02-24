CREATE TABLE IF NOT EXISTS transactions (
    uid BIGSERIAL PRIMARY KEY,
    origin VARCHAR(32) NOT NULL,
    destination VARCHAR(32) NOT NULL,
    amount NUMERIC(18, 2) NOT NULL CHECK (amount >= 0),
    message VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Chave estrangeira para a carteira de origem
    CONSTRAINT transactions_origin_fk FOREIGN KEY (origin) 
        REFERENCES wallets(public_key) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
    
    -- Chave estrangeira para a carteira de destino (ajustado o comportamento ao excluir)
    CONSTRAINT transactions_destination_fk FOREIGN KEY (destination) 
        REFERENCES wallets(public_key) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
);

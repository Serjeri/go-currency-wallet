CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS wallet (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    usd BIGINT DEFAULT 0,
    rub BIGINT DEFAULT 0,
    eur BIGINT DEFAULT 0,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS exchanger (
    id SERIAL PRIMARY KEY,
    usd BIGINT DEFAULT 0,
    rub BIGINT DEFAULT 0,
    eur BIGINT DEFAULT 0,
    usd_rub BIGINT DEFAULT 0,
    usd_eur BIGINT DEFAULT 0,
    eur_rub BIGINT DEFAULT 0,
    eur_usd BIGINT DEFAULT 0,
    rub_eur BIGINT DEFAULT 0,
    rub_usd BIGINT DEFAULT 0,
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_token ON users(token);
CREATE INDEX idx_wallet_user ON wallet(user_id);

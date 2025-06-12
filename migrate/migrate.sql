\c wallet

DO $$
BEGIN
 IF NOT EXISTS (
 SELECT 1
 FROM pg_database
 WHERE datname = 'wallet'
 ) THEN
 RAISE NOTICE 'База данных wallet не существует, создаём...';
 CREATE DATABASE wallet;
 END IF;
END
$$;

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
    rub_usd BIGINT DEFAULT 0
);


CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_wallet_user ON wallet(user_id);


INSERT INTO exchanger (usd, rub, eur, usd_rub, usd_eur, eur_rub, eur_usd, rub_eur, rub_usd)
VALUES (
    100000,
    10000000,
    95000,
    9000,
    95,
    8500,
    105,
    117,
    111
);

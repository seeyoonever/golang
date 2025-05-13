CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT NOT NULL UNIQUE,
    steam_id VARCHAR(64) NOT NULL,
    info TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
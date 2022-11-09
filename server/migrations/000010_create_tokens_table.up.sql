-- Filename: migrations/000010_create_tokens_table.up.sql

CREATE TABLE
    IF NOT EXISTS tokens(
        hash bytea PRIMARY KEY,
        user_id BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
        expiry TIMESTAMP(0)
        WITH
            TIME ZONE NOT NULL,
            scope text NOT NULL
    )
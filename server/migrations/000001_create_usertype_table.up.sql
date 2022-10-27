-- Filename: migrations/000001_create_usertype_table.up.sql

CREATE TABLE
    IF NOT EXISTS usertype(
        id bigserial PRIMARY KEY,
        name text
    );
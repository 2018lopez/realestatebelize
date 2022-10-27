--Filename: migrations/000002_create_district_table.up.sql

CREATE TABLE
    IF NOT EXISTS district(
        id bigserial PRIMARY KEY,
        name text
    );
-- Filname: migriations/000005_create_propertystatus_table.up.sql

CREATE TABLE
    IF NOT EXISTS propertystatus(
        id bigserial PRIMARY KEY,
        name text
    );
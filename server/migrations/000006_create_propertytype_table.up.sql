--Filename : migrations/000006_create_propertytype_table.up.sql

CREATE TABLE
    IF NOT EXISTS propertytype(
        id bigserial PRIMARY KEY,
        name text
    );
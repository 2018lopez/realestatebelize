-- Filename: migrations/000012_add_permissions.up.sql

CREATE TABLE
    IF NOT EXISTS permissions(
        id bigserial PRIMARY KEY,
        code text NOT NULL
    );

-- CREAT A LINKING TABLE THAT links users to permissions

-- many to many  relationship base

CREATE TABLE
    IF NOT EXISTS users_permissions(
        user_id BIGINT not NULL REFERENCES users(id) ON DELETE CASCADE,
        permission_id BIGINT NOT NULL REFERENCES permissions (id) ON DELETE CASCADE,
        PRIMARY KEY(user_id, permission_id)
    );

INSERT INTO permissions(code)
VALUES
('listings:read'), ('listings:write');
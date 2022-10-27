--Filename: migrations/000003_create_users_table.up.sql

CREATE TABLE
    IF NOT EXISTS users(
        id bigserial PRIMARY KEY,
        username text UNIQUE NOT NULL,
        password_hash bytea NOT NULL,
        fullname text NOT NULL,
        email text UNIQUE NOT NULL,
        phone INT NOT NULL,
        profileImageUrl text NOT NULL,
        address text NOT NULL,
        districtId INT NOT NULL,
        userTypeId INT NOT NULL,
        activated BOOL,
        created_at timestamp(0)
        with
            time zone NOT NULL DEFAULT NOW(),
            FOREIGN KEY(userTypeId) REFERENCES usertype(id),
            FOREIGN KEY(districtId) REFERENCES district(id)
    );
-- Filaname - migrations/000011_user_profile_image.up.sql

CREATE TABLE
    IF NOT EXISTS userprofileimage(
        user_id BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
        image_url text NOT NULL
    );
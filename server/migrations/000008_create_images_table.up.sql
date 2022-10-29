--Filename: migrations/000008_create_images_table.up.sql

CREATE TABLE
    IF NOT EXISTS images(
        id bigserial PRIMARY KEY,
        listingId INT,
        imageUrl text NOT NULL,
        FOREIGN KEY(listingId) REFERENCES listing(id)
    );
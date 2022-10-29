--Filename: migrations/000009_create_userproperties_table.up.sql

CREATE TABLE
    IF NOT EXISTS userproperties(
        userId INT NOT NULL,
        listingId INT NOT NULL,
        FOREIGN KEY(listingId) REFERENCES listing(id),
        FOREIGN KEY(userId) REFERENCES users(id)
    );
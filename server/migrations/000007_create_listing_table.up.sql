-- Filename: migrations/000007_create_listing_table.down.sql

CREATE TABLE
    IF NOT EXISTS listing(
        id bigserial PRIMARY KEY,
        propertyTitle text NOT NULL,
        propertyStatusId int,
        propertyTypeId int,
        price decimal NOT NULL,
        description text NOT NULL,
        address text,
        districtId INT,
        googleMapUrl text NOT NULL,
        created_at timestamp(0)
        with
            time zone NOT NULL DEFAULT NOW(),
            FOREIGN KEY(propertyStatusId) REFERENCES propertystatus(id),
            FOREIGN KEY(propertyTypeId) REFERENCES propertytype(id),
            FOREIGN KEY(districtId) REFERENCES district(id)
    );
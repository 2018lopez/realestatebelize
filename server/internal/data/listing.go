// Filename - internal/data/listing.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"realestatebelize.imerlopez.net/internal/validator"
)

type Listing struct {
	ID               int64     `json:"id"`
	PropertyTitle    string    `json:"property_title"`
	PropertyStatusId int64     `json:"property_status_id"`
	PropertyTypeId   int64     `json:"property_type_id"`
	Price            float64   `json:"price"`
	Description      string    `json:"description"`
	Address          string    `json:"address"`
	DistrictId       int64     `json:"district_id"`
	GoogleMapUrl     string    `json:"google_map_url"`
	CreatedAt        time.Time `json:"-"`
}

type Listings struct {
	ID               int64     `json:"id"`
	PropertyTitle    string    `json:"property_title"`
	PropertyStatusId string    `json:"property_status_id"`
	PropertyTypeId   string    `json:"property_type_id"`
	Price            float64   `json:"price"`
	Description      string    `json:"description"`
	Address          string    `json:"address"`
	DistrictId       string    `json:"district_id"`
	GoogleMapUrl     string    `json:"google_map_url"`
	CreatedAt        time.Time `json:"-"`
}

func ValidateListing(v *validator.Validator, listing *Listing) {

	//use the check method to execute our validation
	v.Check(listing.PropertyTitle != "", "property_title", "must be provided")
	v.Check(len(listing.PropertyTitle) >= 20, "propertyt_itle", "must be more than 20 byte long")

	v.Check(listing.PropertyStatusId > 0, "property_status_id", "must be provided")
	v.Check(listing.PropertyTypeId > 0, "property_type_id", "must be provided")

	v.Check(listing.Price >= 0, "price", "must be provided")

	v.Check(listing.Description != "", "description", "must be provided")
	v.Check(len(listing.Description) >= 20, "description", "must be more than 20 byte long")

	v.Check(listing.Address != "", "address", "must be provided")
	v.Check(len(listing.Address) >= 10, "address", "must be more than 10 byte long")

	v.Check(listing.DistrictId > 0, "district_id", "must be provided")

	v.Check(listing.GoogleMapUrl != "", "google_map_url", "must be provided")

}

// Define a ListingModel which wrap a sql.DB connection pool
type ListingModel struct {
	DB *sql.DB
}

// insert() allow us to create a new school
func (m ListingModel) Insert(listing *Listing) error {

	query := `
		INSERT INTO listing(propertytitle,propertystatusid,propertytypeid,price,description,address,districtid,googlemapurl)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	//cleanup to prevent memory leak
	defer cancel()
	//Collect the data fields into a slice
	args := []interface{}{
		listing.PropertyTitle,
		listing.PropertyStatusId,
		listing.PropertyTypeId,
		listing.Price,
		listing.Description,
		listing.Address,
		listing.DistrictId,
		listing.GoogleMapUrl,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&listing.ID, &listing.CreatedAt)

}

// Get () allow us to retrieve a specific school
func (m ListingModel) Get(id int64) (*Listings, error) {

	//Ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	//create query
	query := `

	SELECT l.id , l.propertytitle as title, ps.name as propertystatus, pt.name as propertytype, l.price, l.description, l.address, d.name as district, l.googlemapurl, l.created_at  from listing l inner join propertystatus ps on l.propertystatusid=ps.id
	inner join propertytype pt on l.propertytypeid = pt.id
	inner join district d on l.districtid = d.id
	WHERE l.id = $1
	
	`
	//Declare school variable to hold the return data

	var listing Listings

	//create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	//cleanup to prevent memory leak
	defer cancel()

	//Execute the query using QueryRow()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(

		&listing.ID,
		&listing.PropertyTitle,
		&listing.PropertyStatusId,
		&listing.PropertyTypeId,
		&listing.Price,
		&listing.Description,
		&listing.Address,
		&listing.DistrictId,
		&listing.GoogleMapUrl,
		&listing.CreatedAt,
	)

	if err != nil {
		//check type of err
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	//Success
	return &listing, nil
}

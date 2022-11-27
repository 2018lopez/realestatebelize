// Filename - internal/data/listing.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
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

// listing struct for get by id
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
	Images           []string  `json:"images"`
	Agent            string    `json:"agent"`
	AgentPhone       string    `json:"agent_phone"`
	AgentEmail       string    `json:"agent_email"`
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

func ValidateListings(v *validator.Validator, listing *Listings) {

	//use the check method to execute our validation
	v.Check(listing.PropertyTitle != "", "property_title", "must be provided")
	v.Check(len(listing.PropertyTitle) >= 20, "propertyt_itle", "must be more than 20 byte long")

	v.Check(listing.PropertyStatusId != "", "property_status_id", "must be provided")
	v.Check(listing.PropertyTypeId != "", "property_type_id", "must be provided")

	v.Check(listing.Price >= 0, "price", "must be provided")

	v.Check(listing.Description != "", "description", "must be provided")
	v.Check(len(listing.Description) >= 20, "description", "must be more than 20 byte long")

	v.Check(listing.Address != "", "address", "must be provided")
	v.Check(len(listing.Address) >= 10, "address", "must be more than 10 byte long")

	v.Check(listing.DistrictId != "", "district_id", "must be provided")

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

// Update Listing
func (m ListingModel) Update(listing *Listings) error {

	query := `
	UPDATE listing
	set propertytitle = $1, propertystatusid = (select id from propertystatus where name = $2), propertytypeid = (select id from propertytype where name = $3)
	,price = $4, description = $5, address = $6, districtid = (select id from district where name = $7), googlemapurl = $8
	where id = $9
		RETURNING id
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
		listing.ID,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&listing.ID)

}

// Get () allow us to retrieve a specific listing
func (m ListingModel) Get(id int64) (*Listings, error) {

	//Ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	//create query
	query := `

	SELECT l.id , l.propertytitle as title, ps.name as propertystatus, pt.name as propertytype, l.price, l.description, l.address, d.name as district, l.googlemapurl, i.imageurl,u.fullname, u.phone, u.email, l.created_at  from listing l inner join propertystatus ps on l.propertystatusid=ps.id
	inner join propertytype pt on l.propertytypeid = pt.id
	inner join district d on l.districtid = d.id
	inner join userproperties up on up.listingid = l.id
	inner join users u on u.id = up.userid
	inner join images i on i.listingid = l.id
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
		pq.Array(&listing.Images),
		&listing.Agent,
		&listing.AgentPhone,
		&listing.AgentEmail,
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

// Display all listings
func (m ListingModel) ShowListings(propertytitle string, district string, filters Filters) ([]*Listings, Metadata, error) {

	//create query
	query := fmt.Sprintf(`

	SELECT COUNT(*) OVER(), l.id , l.propertytitle as title, ps.name as propertystatus, pt.name as propertytype, l.price, l.description, l.address, d.name as district, l.googlemapurl, i.imageurl,u.fullname, u.phone, u.email, l.created_at  from listing l inner join propertystatus ps on l.propertystatusid=ps.id
	inner join propertytype pt on l.propertytypeid = pt.id
	inner join district d on l.districtid = d.id
	inner join userproperties up on up.listingid = l.id
	inner join users u on u.id = up.userid
	inner join images i on i.listingid = l.id
	where (to_tsvector('simple', l.propertytitle) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND (to_tsvector('simple', d.name) @@ plainto_tsquery('simple', $2) OR $2 = '')
	ORDER BY %s %s, l.id ASC
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortOrder())

	//create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	//cleanup to prevent memory leak
	defer cancel()

	args := []interface{}{propertytitle, district, filters.limit(), filters.offset()}
	//execute
	rows, err := m.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, err
	}

	//close the result set
	defer rows.Close()

	totalRecords := 0

	//Initialize an empty slice to hold listings data
	listings := []*Listings{}

	//iterate over the rows in the result set

	for rows.Next() {
		var listing Listings
		//scan the values from row into school struct
		err := rows.Scan(
			&totalRecords,
			&listing.ID,
			&listing.PropertyTitle,
			&listing.PropertyStatusId,
			&listing.PropertyTypeId,
			&listing.Price,
			&listing.Description,
			&listing.Address,
			&listing.DistrictId,
			&listing.GoogleMapUrl,
			pq.Array(&listing.Images),
			&listing.Agent,
			&listing.AgentPhone,
			&listing.AgentEmail,
			&listing.CreatedAt,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		//add the listings to our slice
		listings = append(listings, &listing)

	}

	//Check for errors after looping through the result set

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// return slice of listings
	return listings, metadata, nil

}

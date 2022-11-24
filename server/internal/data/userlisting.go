//Filename : internal/data/userlisting.go

package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"realestatebelize.imerlopez.net/internal/validator"
)

type UserListings struct {
	Username  string `json:"username"`
	ListingId int64  `json:"listing_id"`
}

type ListingByUser struct {
	Fullname      string   `json:"fullname"`
	PropertyTitle []string `json:"property_title`
	ListingId     []string `json:"listing_id"`
	Total         int64    `json:"total"`
}

func ValidateUserListings(v *validator.Validator, userlisting *UserListings) {

	v.Check(userlisting.Username != "", "username", "must be provided")
	v.Check(userlisting.ListingId != 0, "listing_id", "must be provided")

}

// create  userlisting model
type UserListingsModel struct {
	DB *sql.DB
}

// assign user a property or listing
func (m UserListingsModel) Insert(userlisting *UserListings) error {
	//create our query
	query :=
		`	
		INSERT INTO userproperties(userid, listingid)
		VALUES((select id from users where username = $1), $2)
		RETURNING listingid


	`

	args := []interface{}{
		userlisting.Username,
		userlisting.ListingId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&userlisting.ListingId)

	if err != nil {
		return err
	}
	return nil
}

//get

// Get () allow us to retrieve a specific listing
func (m UserListingsModel) Get(id int64) (*ListingByUser, error) {

	//Ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	//create query
	query := `

	SELECT  u.fullname,  array_agg(l.propertytitle) as properties, array_agg(l.id) as listingid, count(l.id) as total FROM users u INNER JOIN userproperties up on u.id = up.userid
	INNER JOIN listing l ON l.id = up.listingid where u.id = $1 group by u.fullname
	
	`
	//Declare school variable to hold the return data

	var listing ListingByUser

	//create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	//cleanup to prevent memory leak
	defer cancel()

	//Execute the query using QueryRow()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(

		&listing.Fullname,
		pq.Array(&listing.PropertyTitle),
		pq.Array(&listing.ListingId),
		&listing.Total,
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

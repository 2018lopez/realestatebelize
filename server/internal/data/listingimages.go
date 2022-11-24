//Filename: internal/data/listingimages.go

package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"realestatebelize.imerlopez.net/internal/validator"
)

type ListingImages struct {
	ListingID int64    `json:"listing_id"`
	ImageURl  []string `json:"image_url"`
}

func ValidateListingImages(v *validator.Validator, listingimg *ListingImages) {

	v.Check(listingimg.ListingID != 0, "lsiting_id", "must be provided")

	v.Check(listingimg.ImageURl != nil, "image_url", "must be provided")
	v.Check(len(listingimg.ImageURl) >= 1, "image_url", "must contain at least 1 entry")
}

// create  user model
type ListingImgModel struct {
	DB *sql.DB
}

// create a new user profile image
func (m ListingImgModel) Insert(listingimg *ListingImages) error {
	//create our query
	query :=
		`	
		INSERT INTO images(listingid, imageurl)
		VALUES($1,$2)
		RETURNING listingid
	`

	args := []interface{}{
		listingimg.ListingID,
		pq.Array(listingimg.ImageURl),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&listingimg.ListingID)

	if err != nil {

		return err
	}
	return nil
}

// get id from listing recent create
func (m ListingImgModel) GetByListingId() (int64, error) {

	query := `
	
		SELECT MAX(id) as id FROM listing

	`

	var userpimg UserProfileImage

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query).Scan(

		&userpimg.UserID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, ErrRecordNotFound
		default:
			return 0, err
		}

	}

	return userpimg.UserID, nil
}

// Fileaname : internal/data/models.go

package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("Edit Conflict")
)

// A wrapper for our data models
type Models struct {
	Users   UserModel
	Listing ListingModel
}

// NewModels allow us to create a new models
func NewModels(db *sql.DB) Models {

	return Models{
		Users:   UserModel{DB: db},
		Listing: ListingModel{DB: db},
	}
}

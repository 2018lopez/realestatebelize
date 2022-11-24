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
	Tokens           TokenModel
	Users            UserModel
	UserProfileImage UserProfileImgModel
	Listing          ListingModel
	Permissions      PermissionsModel
	UserListings     UserListingsModel
}

// NewModels allow us to create a new models
func NewModels(db *sql.DB) Models {

	return Models{
		Tokens:           TokenModel{DB: db},
		Users:            UserModel{DB: db},
		UserProfileImage: UserProfileImgModel{DB: db},
		Listing:          ListingModel{DB: db},
		Permissions:      PermissionsModel{DB: db},
		UserListings:     UserListingsModel{DB: db},
	}
}

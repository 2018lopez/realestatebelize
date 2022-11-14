//Filename interna/data/userprofileimage.go

package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"realestatebelize.imerlopez.net/internal/validator"
)

type UserProfileImage struct {
	UserID   int64  `json:"user_id"`
	ImageURl string `json:"image_url"`
}

func ValidateUserProfileImage(v *validator.Validator, userprofileimg *UserProfileImage) {

	v.Check(userprofileimg.ImageURl != "", "image_url", "must be provided")

}

// create  user model
type UserProfileImgModel struct {
	DB *sql.DB
}

// create a new user profile image
func (m UserProfileImgModel) Insert(userprofileimg *UserProfileImage) error {
	//create our query
	query :=
		`	
		INSERT INTO userprofileimage(user_id, image_url)
		VALUES($1,$2)
		RETURNING user_id
	`

	args := []interface{}{
		userprofileimg.UserID,
		userprofileimg.ImageURl,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&userprofileimg.UserID)

	if err != nil {
		switch {
		case err.Error() == `pq :duplicate key values violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

// get user id
func (m UserProfileImgModel) GetByUserId() (int64, error) {

	query := `
	
		SELECT MAX(id) as id FROM users

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

// update user profile image
func (m UserProfileImgModel) Update(userprofileimg *UserProfileImage) error {
	//create our query
	query :=
		`	
		UPDATE userprofileimage
		set image_url = $1
		where user_id = (select id from users where username = $2)
		RETURNING user_id
	`

	args := []interface{}{
		userprofileimg.UserID,
		userprofileimg.ImageURl,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&userprofileimg.UserID)

	if err != nil {
		switch {
		case err.Error() == `pq :duplicate key values violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

// get user id
func (m UserProfileImgModel) GetIdByUsername(username string) (int64, error) {

	query := `
	
		SELECT id from users where username = $1

	`

	var userpimg UserProfileImage

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, username).Scan(

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

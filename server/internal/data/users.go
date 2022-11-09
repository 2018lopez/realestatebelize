// Filename - internal/data/users.go
package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"realestatebelize.imerlopez.net/internal/validator"
)

var (
	ErrDuplicateEmail = errors.New("Duplicate email")
)

type User struct {
	ID              int64     `json:"id"`
	Username        string    `json:"username"`
	Password        password  `json:"-"`
	Fullname        string    `json:"fullname"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	Address         string    `json:"address"`
	DistrictId      int64     `json:"district_id"`
	ProfileImageUrl string    `json:"profile_image_url"`
	UserTypeId      int64     `json:"user_type_id"`
	Activated       bool      `json:"activated"`
	CreatedAt       time.Time `json:"created_at"`
}

// create a customer password type
type password struct {
	plaintext *string
	hash      []byte
}

// The set method to stores the hash of plaintext password
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

// The matches() method checks of the supplied password is correct
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

//Validate the Client request

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be more than 8 Characters")
	v.Check(len(password) <= 72, "password", "must not be more than 72 Characters")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Username != "", "username", "must be provided")
	v.Check(len(user.Username) <= 200, "username", "must not be more than 200 bye Characters")

	v.Check(user.Fullname != "", "fullname", "must be provided")
	v.Check(len(user.Fullname) <= 500, "fullname", "must not be more than 500 bye Characters")

	v.Check(user.Phone != "", "phone", "must be provided")
	v.Check(validator.Matches(user.Phone, validator.PhoneRX), "phone", "must be a valid phone number")

	v.Check(user.Address != "", "address", "must be provided")
	v.Check(len(user.Address) <= 500, "address", "must not be more than 500 bytes long")

	v.Check(user.DistrictId != 0, "district_id", "must be provided")
	v.Check(user.UserTypeId != 0, "user_type_id", "must be provided")

	//validate Email

	ValidateEmail(v, user.Email)

	//validate Password
	if user.Password.plaintext != nil {

		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	//Ensure a hash of the password was created

	if user.Password.hash == nil {
		panic("Missing password hash for the user")
	}

}

// create  user model
type UserModel struct {
	DB *sql.DB
}

// create a new user
func (m UserModel) Insert(user *User) error {
	//create our query
	query :=
		`	
		INSERT INTO users(username, password_hash, fullname, email,phone, profileimageurl,address, districtid,usertypeid,activated)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id, created_at
	`

	args := []interface{}{
		user.Username,
		user.Password.hash,
		user.Fullname,
		user.Email,
		user.Phone,
		user.ProfileImageUrl,
		user.Address,
		user.DistrictId,
		user.UserTypeId,
		user.Activated,
	}

	fmt.Println("Back", args)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)

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

// The client can update their information
func (m UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET username = $1, password_hash = $2, fullname = $3, email = $4, phone = $5,
		profileimageurl = $6, address = $7, districtid = $8, usertypeid = $9 , activated = $10
		WHERE id = $11
		RETURNING id
	`
	args := []interface{}{
		user.Username,
		user.Password.hash,
		user.Fullname,
		user.Email,
		user.Phone,
		user.ProfileImageUrl,
		user.Address,
		user.DistrictId,
		user.UserTypeId,
		user.Activated,
		user.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID)
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

func (m UserModel) GetForToken(tokenScope, tokenPlainText string) (*User, error) {

	tokenHash := sha256.Sum256([]byte(tokenPlainText))
	//setup query
	query := `

		SELECT users.id, users.username, users.password_hash, users.fullname, users.email, users.phone, users.profileimageurl, users.address, users.districtid, users.usertypeid, users.activated, users.created_at
		FROM users
		INNER JOIN tokens
		on users.id = tokens.user_id
		WHERE tokens.hash = $1 AND tokens.scope = $2
		AND tokens.expiry > $3
	`
	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Password.hash,
		&user.Fullname,
		&user.Email,
		&user.Phone,
		&user.ProfileImageUrl,
		&user.Address,
		&user.DistrictId,
		&user.UserTypeId,
		&user.Activated,
		&user.CreatedAt,
	)

	if err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil

}

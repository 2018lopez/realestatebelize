// Filename - internal/data/users.go
package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"realestatebelize.imerlopez.net/internal/validator"
)

var (
	ErrDuplicateEmail = errors.New("Duplicate email")
	AnonymousUser     = &User{}
)

type User struct {
	ID         int64    `json:"id"`
	Username   string   `json:"username"`
	Password   password `json:"-"`
	Fullname   string   `json:"fullname"`
	Email      string   `json:"email"`
	Phone      string   `json:"phone"`
	Address    string   `json:"address"`
	DistrictId int64    `json:"district_id"`

	UserTypeId int64     `json:"user_type_id"`
	Activated  bool      `json:"activated"`
	CreatedAt  time.Time `json:"created_at"`
}

//UserListing struct use for get by id

type UserListing struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Password     password  `json:"-"`
	Fullname     string    `json:"fullname"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	DistrictId   string    `json:"district_id"`
	UserTypeId   string    `json:"user_type_id"`
	Activated    bool      `json:"activated"`
	ProfileImage string    `json:"profile_image"`
	CreatedAt    time.Time `json:"created_at"`
}

// User Struct for password reset
type UserResetPassword struct {
	Username string   `json:"username"`
	Password password `json:"-"`
}

// create a customer password type
type password struct {
	plaintext *string
	hash      []byte
}

// check if a user is anonymous
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
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

// Validate for user struct
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

// Validate the userlisting struct
func ValidateUserListing(v *validator.Validator, user *UserListing) {
	v.Check(user.Username != "", "username", "must be provided")
	v.Check(len(user.Username) <= 200, "username", "must not be more than 200 bye Characters")

	v.Check(user.Fullname != "", "fullname", "must be provided")
	v.Check(len(user.Fullname) <= 500, "fullname", "must not be more than 500 bye Characters")

	v.Check(user.Phone != "", "phone", "must be provided")
	v.Check(validator.Matches(user.Phone, validator.PhoneRX), "phone", "must be a valid phone number")

	v.Check(user.Address != "", "address", "must be provided")
	v.Check(len(user.Address) <= 500, "address", "must not be more than 500 bytes long")

	v.Check(user.DistrictId != "", "district_id", "must be provided")
	v.Check(user.UserTypeId != "", "user_type_id", "must be provided")

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

// Validate the userlisting struct
func ValidateUserReset(v *validator.Validator, user *UserResetPassword) {
	v.Check(user.Username != "", "username", "must be provided")
	v.Check(len(user.Username) <= 200, "username", "must not be more than 200 bye Characters")

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
		INSERT INTO users(username, password_hash, fullname, email,phone, address, districtid,usertypeid,activated)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at
	`

	args := []interface{}{
		user.Username,
		user.Password.hash,
		user.Fullname,
		user.Email,
		user.Phone,
		user.Address,
		user.DistrictId,
		user.UserTypeId,
		user.Activated,
	}

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
		 address = $6, districtid = $7, usertypeid = $8 , activated = $9
		WHERE id = $10
		RETURNING id
	`
	args := []interface{}{
		user.Username,
		user.Password.hash,
		user.Fullname,
		user.Email,
		user.Phone,
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

// Update userListing - its different since every entries is string type so additional internal need to place for a successful update
// The client can update their information
func (m UserModel) UpdateUser(user *UserListing) error {
	query := `
		UPDATE users
		SET username = $1, fullname = $2, email = $3, phone = $4,
		 address = $5, districtid = (select id from district where name = $6), usertypeid = (select id from usertype where name = $7) , activated = $8
		WHERE id = $9
		RETURNING id
	`
	args := []interface{}{
		user.Username,
		user.Fullname,
		user.Email,
		user.Phone,
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

// User - Reset password
func (m UserModel) ResetPassword(user *UserResetPassword) error {
	query := `
		UPDATE users
		SET password_hash = $1
		WHERE username = $2
		RETURNING username
	`
	args := []interface{}{

		user.Password.hash,
		user.Username,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Username)
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

		SELECT users.id, users.username, users.password_hash, users.fullname, users.email, users.phone,  users.address, users.districtid, users.usertypeid, users.activated, users.created_at
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

// get user based on their username
func (m UserModel) GetByUsername(username string) (*User, error) {

	query := `
	
		SELECT id, username, password_hash, fullname, email,phone, address, districtid,usertypeid,activated, created_at
		FROM users
		WHERE username = $1
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, username).Scan(

		&user.ID,
		&user.Username,
		&user.Password.hash,
		&user.Fullname,
		&user.Email,
		&user.Phone,

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

// Get () allow us to retrieve a specific listing
func (m UserModel) Get(id int64) (*UserListing, error) {

	//Ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	//create query
	query := `

	SELECT u.id, u.username, u.password_hash,u.fullname, u.email, u.phone, u.address, d.name as district, ut.name as usertype,
		u.activated, img.image_url, u.created_at
		FROM users u inner join userprofileimage img
		on u.id = img.user_id
		inner join district d 
		on d.id = u.districtid
		inner join usertype ut
		on ut.id = u.usertypeid
		where u.id =  $1
	
	`
	//Declare school variable to hold the return data

	var userlisting UserListing

	//create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	//cleanup to prevent memory leak
	defer cancel()

	//Execute the query using QueryRow()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(

		&userlisting.ID,
		&userlisting.Username,
		&userlisting.Password.hash,
		&userlisting.Fullname,
		&userlisting.Email,
		&userlisting.Phone,
		&userlisting.Address,
		&userlisting.DistrictId,
		&userlisting.UserTypeId,
		&userlisting.Activated,
		&userlisting.ProfileImage,
		&userlisting.CreatedAt,
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
	return &userlisting, nil
}

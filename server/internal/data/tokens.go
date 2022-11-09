//Filename: internal/data/tokens.go

package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"realestatebelize.imerlopez.net/internal/validator"
)

//token categories/scopes

const (
	ScopeActivation = "activation"
)

//Define token type

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

//generateToken() function returns a token

func generateTokenT(userID int64, ttl time.Duration, scope string) (*Token, error) {

	//create token
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	//create byte slices to hold random values & fill it with values
	//from CSPRNG
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	//encode the byte slice to a base32 encoding string
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	//Hash the string token
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil

}

// check that the plaintext token is 26 bytes
func ValidateTokenPlainText(v *validator.Validator, tokenPlainText string) {

	v.Check(tokenPlainText != "", "token", "must be 26 bytes long")
	v.Check(len(tokenPlainText) == 26, "token", "must be 26 bytes long")
}

// define token model
type TokenModel struct {
	DB *sql.DB
}

// create an insert a token into token table
func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {

	token, err := generateTokenT(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

//insert entry to tokens table

func (m TokenModel) Insert(token *Token) error {

	query := `

		INSERT INTO tokens( hash, user_id, expiry, scope)
		VALUES($1,$2,$3,$4)
	
	`

	args := []interface{}{
		token.Hash,
		token.UserID,
		token.Expiry,
		token.Scope,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err

}

// Delete token
func (m TokenModel) DeleteAllForUsers(scope string, userID int64) error {

	query := `
		DELETE FROM tokens WHERE scope = $1 and user_id = $2

	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err

}

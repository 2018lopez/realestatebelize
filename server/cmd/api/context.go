//Filename : cmd/api/context.go

package main

import (
	"context"
	"net/http"

	"realestatebelize.imerlopez.net/internal/data"
)

// Define a custom contextKey type
type contextKey string

// make user a key
const userContextKey = contextKey("user")

// method to add user to the context
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {

	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

//Retreive the User Struct

func (app *application) contextGetUser(r *http.Request) *data.User {

	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}

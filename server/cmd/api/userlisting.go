//Filename: cmd/api/userlisitng.go

package main

import (
	"errors"
	"net/http"

	"realestatebelize.imerlopez.net/internal/data"
	"realestatebelize.imerlopez.net/internal/validator"
)

func (app *application) addUserListingHandler(w http.ResponseWriter, r *http.Request) {
	//our target decode distination
	var input struct {
		Username  string `json:"username"`
		ListingId int64  `json:"listing_id"`
	}

	//initialize the new json decoder instance

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//copy the values from the input struct to a new listing struct
	userlisting := &data.UserListings{
		Username:  input.Username,
		ListingId: input.ListingId,
	}

	//Initialize a new Validator instance
	v := validator.New()

	//check the map to determine if there were any validation errors
	if data.ValidateUserListings(v, userlisting); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//add listing to agent

	err = app.models.UserListings.Insert(userlisting)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	//write json response with 201

	err = app.writeJSON(w, http.StatusCreated, envelope{"User_Listin": userlisting}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// get properties list by agent id
func (app *application) getListingByAgentdHandler(w http.ResponseWriter, r *http.Request) {

	//get id from param
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//fetch the specific schools

	users, err := app.models.UserListings.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return

	}

	//write data return by get
	err = app.writeJSON(w, http.StatusOK, envelope{"Listing_User": users}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

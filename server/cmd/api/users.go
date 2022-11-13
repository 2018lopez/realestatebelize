//Filename: cmd/api/users.go

package main

import (
	"errors"
	"net/http"
	"time"

	"realestatebelize.imerlopez.net/internal/data"
	"realestatebelize.imerlopez.net/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	//Hold data from the request body

	var input struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		Fullname   string `json:"fullname"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Address    string `json:"address"`
		DistrictId int64  `json:"district_id"`
		UserTypeId int64  `json:"user_type_id"`
	}

	//Parse the request body into the anonymous struct
	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Username:   input.Username,
		Fullname:   input.Fullname,
		Email:      input.Email,
		Phone:      input.Phone,
		Address:    input.Address,
		DistrictId: input.DistrictId,
		UserTypeId: input.UserTypeId,
		Activated:  false,
	}

	//Generate a password Hash
	err = user.Password.Set(input.Password)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//Perform Validation
	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//insert data to database

	err = app.models.Users.Insert(user)
	if err != nil {

		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	//genrate a token for the newly-create user
	token, err := app.models.Tokens.New(user.ID, 1*24*time.Hour, data.ScopeActivation)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {

		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		//Send the mail to the new user
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			// log errors
			app.logger.PrintError(err, nil)
		}

	})
	//Write a 202 Accepted Status
	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Update User
func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {

	//get id from user
	id, err := app.readIdParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//fetch the orignal record from the database
	user, err := app.models.Users.Get(id)

	if err != nil {

		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	//Create an input Struct to hold data read in from client

	var input struct {
		Username   *string `json:"username"`
		Fullname   *string `json:"fullname"`
		Email      *string `json:"email"`
		Phone      *string `json:"phone"`
		Address    *string `json:"address"`
		DistrictId *string `json:"district_id"`
		UserTypeId *string `json:"user_type_id"`
		Activated  bool    `json:"activated"`
	}
	//intialize new json.decoder instance

	err = app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// //Check the update
	if input.Username != nil {

		user.Username = *input.Username
	}
	if input.Fullname != nil {

		user.Fullname = *input.Fullname
	}

	if input.Email != nil {

		user.Email = *input.Email
	}

	if input.Phone != nil {

		user.Phone = *input.Phone
	}

	if input.Address != nil {

		user.Address = *input.Address
	}

	if input.DistrictId != nil {
		user.DistrictId = *input.DistrictId
	}

	if input.UserTypeId != nil {
		user.UserTypeId = *input.UserTypeId
	}

	//Initalize a new Validator
	v := validator.New()

	//check the map to determine if there were any validation errors

	if data.ValidateUserListing(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return

	}

	err = app.models.Users.UpdateUser(user)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	//wreite data return by the update
	err = app.writeJSON(w, http.StatusOK, envelope{"users": user}, nil)

}

func (app *application) activatedUserHandler(w http.ResponseWriter, r *http.Request) {

	//parse the plaintext activation token
	var input struct {
		TokenPlainText string `json:"token"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//Perform Validation

	v := validator.New()

	if data.ValidateTokenPlainText(v, input.TokenPlainText); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//get the user info of the provide token or give
	//client feedback regarding invalid token

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlainText)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	//update the user status

	user.Activated = true

	//save the update user's record on the database

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//delete the user's token that was used for activation
	err = app.models.Tokens.DeleteAllForUsers(data.ScopeActivation, user.ID)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//Send JSON response with the update detail
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// Show user for get by id

func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {

	//get id from param
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//fetch the specific schools

	users, err := app.models.Users.Get(id)
	//handle errors
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
	err = app.writeJSON(w, http.StatusOK, envelope{"user": users}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

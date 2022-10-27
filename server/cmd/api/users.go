//Filename: cmd/api/users.go

package main

import (
	"errors"
	"net/http"

	"realestatebelize.imerlopez.net/internal/data"
	"realestatebelize.imerlopez.net/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	//Hold data from the request body

	var input struct {
		Username        string `json:"username"`
		Password        string `json:"password"`
		Fullname        string `json:"fullname"`
		Email           string `json:"email"`
		Phone           string `json:"phone"`
		Address         string `json:"address"`
		DistrictId      int64  `json:"district_id"`
		ProfileImageUrl string `json:"profile_image_url"`
		UserTypeId      int64  `json:"user_type_id"`
	}

	//Parse the request body into the anonymous struct
	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Username:        input.Username,
		Fullname:        input.Fullname,
		Email:           input.Email,
		Phone:           input.Phone,
		Address:         input.Address,
		DistrictId:      input.DistrictId,
		ProfileImageUrl: input.ProfileImageUrl,
		UserTypeId:      input.UserTypeId,
		Activated:       false,
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

	//Write a 201 status

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

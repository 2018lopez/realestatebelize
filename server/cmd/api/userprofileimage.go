//Filename: cmd/api/userprofileimage.go

package main

import (
	"errors"
	"net/http"

	"realestatebelize.imerlopez.net/internal/data"
	"realestatebelize.imerlopez.net/internal/validator"
)

func (app *application) uploadUserImageHandler(w http.ResponseWriter, r *http.Request) {

	// var input struct {
	// 	ImageUrl string `json:"image_url"`
	// }

	//get the id for last user account created

	user, err := app.models.UserProfileImage.GetByUserId()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	imagePath, err := app.uploadFiles(r)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	userimg := &data.UserProfileImage{
		UserID:   user,
		ImageURl: imagePath,
	}

	//Perform Validation
	v := validator.New()

	if data.ValidateUserProfileImage(v, userimg); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.UserProfileImage.Insert(userimg)
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

	//Send JSON response with the update detail
	err = app.writeJSON(w, http.StatusOK, envelope{"user_profile_image": userimg}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

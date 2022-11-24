//Filename: cmd/api/listingimages.go

package main

import (
	"fmt"
	"net/http"

	"realestatebelize.imerlopez.net/internal/data"
	"realestatebelize.imerlopez.net/internal/validator"
)

// upload listing images
func (app *application) uploadListingImageHandler(w http.ResponseWriter, r *http.Request) {

	//get the id for last listing created

	listing, err := app.models.ListingImages.GetByListingId()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	imagePath, err := app.uploadImages(r)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	fmt.Println(imagePath)

	listingimg := &data.ListingImages{
		ListingID: listing,
		ImageURl:  imagePath,
	}

	//Perform Validation
	v := validator.New()

	if data.ValidateListingImages(v, listingimg); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.ListingImages.Insert(listingimg)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	//Send JSON response with the update detail
	err = app.writeJSON(w, http.StatusOK, envelope{"listing_images": listingimg}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

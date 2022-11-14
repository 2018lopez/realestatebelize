//Filename: cmd/api/listing.go

package main

import (
	"errors"
	"fmt"
	"net/http"

	"realestatebelize.imerlopez.net/internal/data"
	"realestatebelize.imerlopez.net/internal/validator"
)

func (app *application) createListingHandler(w http.ResponseWriter, r *http.Request) {
	//our target decode distination
	var input struct {
		PropertyTitle    string  `json:"property_title"`
		PropertyStatusId int64   `json:"property_status_id"`
		PropertyTypeId   int64   `json:"property_type_id"`
		Price            float64 `json:"price"`
		Description      string  `json:"description"`
		Address          string  `json:"address"`
		DistrictId       int64   `json:"district_id"`
		GoogleMapUrl     string  `json:"google_map_url"`
	}

	//initialize the new json decoder instance

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//copy the values from the input struct to a new listing struct
	listing := &data.Listing{
		PropertyTitle:    input.PropertyTitle,
		PropertyStatusId: input.PropertyStatusId,
		PropertyTypeId:   input.PropertyTypeId,
		Price:            input.Price,
		Description:      input.Description,
		Address:          input.Address,
		DistrictId:       input.DistrictId,
		GoogleMapUrl:     input.GoogleMapUrl,
	}

	//Initialize a new Validator instance
	v := validator.New()

	//check the map to determine if there were any validation errors
	if data.ValidateListing(v, listing); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//create listing

	err = app.models.Listing.Insert(listing)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	//create a location header for newly resource : listing
	headers := make(http.Header)
	headers.Set("Locations", fmt.Sprintf("/v1/listings/%d", listing.ID))

	//write json response with 201

	err = app.writeJSON(w, http.StatusCreated, envelope{"listing": listing}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Show listings for get by id

func (app *application) showListingHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//fetch the specific schools

	listing, err := app.models.Listing.Get(id)
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
	err = app.writeJSON(w, http.StatusOK, envelope{"listing": listing}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// Show listings for get by id

func (app *application) updateListingHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//fetch the specific schools

	listing, err := app.models.Listing.Get(id)
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

	//Create an input Struct to hold data read in from client

	var input struct {
		PropertyTitle    *string  `json:"property_title"`
		PropertyStatusId *string  `json:"property_status_id"`
		PropertyTypeId   *string  `json:"property_type_id"`
		Price            *float64 `json:"price"`
		Description      *string  `json:"description"`
		Address          *string  `json:"address"`
		DistrictId       *string  `json:"district_id"`
		GoogleMapUrl     *string  `json:"google_map_url"`
	}
	//intialize new json.decoder instance

	err = app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// //Check the update
	if input.PropertyTitle != nil {

		listing.PropertyTitle = *input.PropertyTitle
	}
	if input.PropertyStatusId != nil {

		listing.PropertyStatusId = *input.PropertyStatusId
	}

	if input.PropertyTypeId != nil {

		listing.PropertyTypeId = *input.PropertyTypeId
	}

	if input.Price != nil {

		listing.Price = *input.Price
	}

	if input.Description != nil {

		listing.Description = *input.Description
	}

	if input.Address != nil {
		listing.Address = *input.Address
	}

	if input.DistrictId != nil {
		listing.DistrictId = *input.DistrictId
	}

	if input.GoogleMapUrl != nil {
		listing.GoogleMapUrl = *input.GoogleMapUrl
	}

	//Initalize a new Validator
	v := validator.New()

	//check the map to determine if there were any validation errors

	if data.ValidateListings(v, listing); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return

	}

	err = app.models.Listing.Update(listing)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}
	//write data return by get
	err = app.writeJSON(w, http.StatusOK, envelope{"listing": listing}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

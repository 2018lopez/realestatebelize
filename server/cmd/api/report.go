// Filaname: cmd/api/report.go
package main

import (
	"net/http"
)

//getTopAgentsHandler allow client to see a top agents

func (app *application) getTopAgentsHandler(w http.ResponseWriter, r *http.Request) {

	//get a listing of the tog agents
	agents, err := app.models.TopAgents.GetTopAgents()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//send a json response
	err = app.writeJSON(w, http.StatusOK, envelope{"top_agents": agents}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

//getListingStatusHandler allow client to see a Total Listing Available or Sold/leased

func (app *application) getListingStatusHandler(w http.ResponseWriter, r *http.Request) {

	//get a listing of the properties
	listing, err := app.models.ListingsStatus.GetListingStatus()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//send a json response
	err = app.writeJSON(w, http.StatusOK, envelope{"listing_status": listing}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

//getTotalSalesHandler allow client to see a TotalSales of properties sold/leased

func (app *application) getTotalSalesHandler(w http.ResponseWriter, r *http.Request) {

	//get a listing of the properties
	sales, err := app.models.TotalSales.GetTotalSales()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//send a json response
	err = app.writeJSON(w, http.StatusOK, envelope{"sales": sales}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

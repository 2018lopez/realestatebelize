//Filename cmd/api/routes

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	//create router
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// File Server Route
	router.ServeFiles("/uploads/*filepath", http.Dir("uploads"))

	//Users routes
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.getUserByIdHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activatedUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/image", app.uploadUserImageHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/image/update", app.updateUserImageHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/updated/:id", app.updateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/resetpassword", app.resetPasswordHandler)
	//End User Routes

	//Listing Routes
	router.HandlerFunc(http.MethodPost, "/v1/listings", app.requirePermission("listings:write", app.createListingHandler))
	router.HandlerFunc(http.MethodPost, "/v1/listings/images", app.uploadListingImageHandler)
	router.HandlerFunc(http.MethodGet, "/v1/listings", app.showAllListingHandler)
	router.HandlerFunc(http.MethodGet, "/v1/listings/:id", app.requirePermission("listings:read", app.showListingHandler))
	router.HandlerFunc(http.MethodPut, "/v1/listings/update/:id", app.updateListingHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/listings", app.addUserListingHandler)
	router.HandlerFunc(http.MethodGet, "/v1/agent/listings/:id", app.getListingByAgentdHandler)
	//End of Listing Routes

	//Report Routes
	router.HandlerFunc(http.MethodGet, "/v1/report/agents", app.getTopAgentsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/report/listings", app.getListingStatusHandler)
	router.HandlerFunc(http.MethodGet, "/v1/report/total-sales", app.getTotalSalesHandler)

	//Currency Rate Route - Third Party API
	router.HandlerFunc(http.MethodGet, "/v1/currencyrate/:id", app.currencyRate)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))

}

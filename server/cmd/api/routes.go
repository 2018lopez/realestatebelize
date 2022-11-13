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
	router.ServeFiles("/uploads/*filepath", http.Dir("uploads"))

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.getUserByIdHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activatedUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/image", app.uploadUserImageHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/listings", app.createListingHandler)
	router.HandlerFunc(http.MethodGet, "/v1/listings/:id", app.showListingHandler)

	return app.recoverPanic(app.rateLimit(router))

}

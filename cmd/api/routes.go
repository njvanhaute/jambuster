package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/tunes", app.requireActivatedUser(app.listTunesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/tunes", app.requireActivatedUser(app.createTuneHandler))

	router.HandlerFunc(http.MethodGet, "/v1/tunes/:id", app.requireActivatedUser(app.showTuneHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/tunes/:id", app.requireActivatedUser(app.updateTuneHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/tunes/:id", app.requireActivatedUser(app.deleteTuneHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}

package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/tunes", app.requirePermission("tunes:read", app.listTunesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/tunes", app.requirePermission("tunes:write", app.createTuneHandler))

	router.HandlerFunc(http.MethodGet, "/v1/tunes/:id", app.requirePermission("tunes:read", app.showTuneHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/tunes/:id", app.requirePermission("tunes:write", app.updateTuneHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/tunes/:id", app.requirePermission("tunes:write", app.deleteTuneHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	router.HandlerFunc(http.MethodPut, "/v1/users/activate", app.activateUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/password", app.updateUserPasswordHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}

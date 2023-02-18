package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.ServeFiles("/templates/*filepath", http.Dir("./templates/"))

	router.HandlerFunc(http.MethodGet, "/", app.homeHandler)

	router.HandlerFunc(http.MethodPost, "/welcome", app.welcomeHandler)
	router.HandlerFunc(http.MethodPost, "/myacc", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodGet, "/activated/:token", app.activateUserHandler)
	router.HandlerFunc(http.MethodGet, "/myacc/tasks", app.requireActivatedUser(app.activateUserHandler))
	router.HandlerFunc(http.MethodGet, "/myacc/upcoming", app.requireActivatedUser(app.activateUserHandler))

	// Return the httprouter instance.
	return app.authenticate(router)
}

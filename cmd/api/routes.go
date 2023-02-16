package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	router.ServeFiles("/templates/*filepath", http.Dir("./templates/"))
	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method. Note that http.MethodGet and
	// http.MethodPost are constants which equate to the strings "GET" and "POST"
	// respectively.
	router.HandlerFunc(http.MethodGet, "/", app.homeHandler)

	router.HandlerFunc(http.MethodPost, "/welcome", app.welcomeHandler)
	router.HandlerFunc(http.MethodPost, "/tasks", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodGet, "/activated/:token", app.activateUserHandler)

	// Return the httprouter instance.
	return router
}

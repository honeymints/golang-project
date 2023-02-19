package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	router.ServeFiles("/templates/*filepath", http.Dir("./templates/"))

	router.HandlerFunc(http.MethodGet, "/", app.homeHandler)

	router.HandlerFunc(http.MethodPost, "/welcome", app.welcomeHandler)
	router.HandlerFunc(http.MethodPost, "/myacc", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodGet, "/activated/:token", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/myacc/tasks", app.createAuthenticationTokenHandler)

	// Return the httprouter instance.
	return router
}

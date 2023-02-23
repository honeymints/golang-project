package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	router.ServeFiles("/templates/*filepath", http.Dir("./templates/"))

	router.HandlerFunc(http.MethodGet, "/login-register", app.loginregisterHandler)
	router.HandlerFunc(http.MethodGet, "/", app.homeHandler)

	router.HandlerFunc(http.MethodPost, "/welcome", app.welcomeHandler)
	router.HandlerFunc(http.MethodPost, "/myacc", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodGet, "/activated/:token", app.activateUserHandler)
	router.HandlerFunc(http.MethodGet, "/myacc/today", app.requireAuthenticatedUser(app.showlistHandler))

	router.HandlerFunc(http.MethodPost, "/myacc/today", app.createlistHandler)
	router.HandlerFunc(http.MethodPost, "/myacc/logout", app.logoutHandler)
	router.HandlerFunc(http.MethodGet, "/myacc/settings", app.requireAuthenticatedUser(app.showlistHandler))

	router.HandlerFunc(http.MethodDeLete, "/myacc/delete", app.deleteAccHandler)
	router.HandlerFunc(http.MethodPut, "/myacc/update", app.UpdateUserHandler)
	//router.HandlerFunc()

	// Return the httprouter instance.ss
	return app.authenticate(router)
}

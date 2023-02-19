package main

import (
	"net/http"
	"time"

	"todolist.net/internal/data"
)

func (app *application) createlistHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user information from the access token.
	user, err := app.models.Users.GetForToken(data.ScopeAuthentication, r.Header.Get("Authorization"))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Create a new to-do list associated with the user.
	list := &data.Lists{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		CreatedAt:   time.Now(),
		User_ID:     user.ID,
	}

	// Insert the list data into the database.
	err = app.models.Lists.Insert(list)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	// Return the newly created to-do list.
	err = app.writeJSON(w, http.StatusAccepted, envelope{"list": list}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

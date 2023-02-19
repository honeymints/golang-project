package main

import (
	"net/http"
	"text/template"
)

func (app *application) createlistHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("templates/tasks.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	dat := r.FormValue("title")
	tpl.ExecuteTemplate(w, "tasks.html", dat)
	/*
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
	*/
}

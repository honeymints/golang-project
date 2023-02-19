package main

import (
	"net/http"
	"text/template"
	"time"

	"todolist.net/internal/data"
)

func (app *application) createlistHandler(w http.ResponseWriter, r *http.Request) {

	list := &data.Lists{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		CreatedAt:   time.Now(),
	}
	//v := validator.New()
	// Insert the list data into the database.
	err := app.models.Lists.Insert(list)
	if err != nil {
		switch {

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// After the user record has been created in the database, generate a new activation
	// token for the user.

	err = app.writeJSON(w, http.StatusAccepted, envelope{"list": list}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	} else {
		tpl, err := template.ParseFiles("templates/tasks.html")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		dat := r.FormValue("txt")

		tpl.ExecuteTemplate(w, "tasks.html", dat)
	}

}

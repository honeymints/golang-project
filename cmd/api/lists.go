package main

import (
	"net/http"
	"text/template"
	"time"

	"todolist.net/internal/data"
)

func (app *application) createlistHandler(w http.ResponseWriter, r *http.Request) {

	cookieToken, _ := r.Cookie("token")
	// Retrieve the user information from the access token.
	user, err := app.models.Users.GetByHash(cookieToken.Value)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Create a new to-do list associated with the user.
	list := &data.Lists{
		Title:     r.FormValue("title"),
		CreatedAt: time.Now(),
		User_ID:   user.ID,
	}

	// Insert the list data into the database.
	err = app.models.Lists.Insert(list)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}
	http.Redirect(w, r, "/myacc/today", http.StatusSeeOther)

	/* lists, err := app.models.Lists.GetByUser(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Return the newly created to-do list.
	err = app.writeJSON(w, http.StatusAccepted, envelope{"list": list}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	tpl, err := template.ParseFiles("templates/tasks.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tpl.ExecuteTemplate(w, "tasks.html", lists) */

}
func (app *application) showlistHandler(w http.ResponseWriter, r *http.Request) {
	cookieToken, _ := r.Cookie("token")
	// Retrieve the user information from the access token.
	user, err := app.models.Users.GetByHash(cookieToken.Value)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	lists, err := app.models.Lists.GetByUser(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Return the list items as JSON.
	err = app.writeJSON(w, http.StatusOK, envelope{"lists": lists}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	tpl, err := template.ParseFiles("templates/tasks.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tpl.ExecuteTemplate(w, "tasks.html", lists)
}

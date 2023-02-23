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

	err = app.writeJSON(w, http.StatusAccepted, envelope{"list": list}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
func (app *application) showlistHandler(w http.ResponseWriter, r *http.Request) {
	cookieToken, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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

	type Data struct {
		Date string
		List []*data.Lists
		User *data.User
	}
	now := time.Now()
	dateStr := now.Format("2006-01-02")
	data := Data{Date: dateStr, List: lists, User: user}
	tpl := template.Must(template.ParseFiles("templates/tasks-2.html"))
	tpl.Execute(w, data)
}

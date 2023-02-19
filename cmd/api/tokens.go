package main

import (
	"errors"
	"net/http"
	"text/template"
	"time"

	"todolist.net/internal/data"
	"todolist.net/internal/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("pswd")
	// Validate the email and password provided by the client.
	v := validator.New()
	data.ValidateEmail(v, email)
	data.ValidatePasswordPlaintext(v, password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Lookup the user record based on the email address. If no matching user was
	// found, then we call the app.invalidCredentialsResponse() helper to send a 401
	// Unauthorized response to the client (we will create this helper in a moment).
	user, err := app.models.Users.GetByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	if !user.Activated {
		app.notActivatedUser(w, r)
	}
	// Check if the provided password matches the actual password for the user.
	match, err := user.Password.Matches(password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// If the passwords don't match, then we call the app.invalidCredentialsResponse()
	// helper again and return.
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}
	// Otherwise, if the password is correct, we generate a new token with a 24-hour
	// expiry time and the scope 'authentication'.
	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Encode the token to JSON and send it in the response along with a 201 Created
	// status code.
	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	http.Redirect(w, r, "/myacc/tasks/", http.StatusSeeOther)

	tpl, err := template.ParseFiles("templates/tasks.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	dat := user.Name

	tpl.ExecuteTemplate(w, "tasks.html", dat)

}

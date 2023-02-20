package main

import (
	"errors"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/julienschmidt/httprouter"
	"todolist.net/internal/data"
	"todolist.net/internal/validator"
)

// Add a createMovieHandler for the "POST /v1/movies" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) loginregisterHandler(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "templates/Untitled-2.html")

}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "templates/main.html")

}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the plaintext activation token from the request body.
	params := httprouter.ParamsFromContext(r.Context())
	tokens := params.ByName("token")

	// Validate the plaintext token provided by the client.
	v := validator.New()
	if data.ValidateTokenPlaintext(v, tokens); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Retrieve the details of the user associated with the token using the
	// GetForToken() method (which we will create in a minute). If no matching record
	// is found, then we let the client know that the token they provided is not valid.
	user, err := app.models.Users.GetForToken(data.ScopeActivation, tokens)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Update the user's activation status.
	user.Activated = true
	// Save the updated user record in our database, checking for any edit conflicts in
	// the same way that we did for our movie records.
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// If everything went successfully, then we delete all activation tokens for the
	// user.
	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send the updated user details to the client in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	} else {
		tpl, err := template.ParseFiles("templates/activated.html")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		dat := user.Name

		tpl.ExecuteTemplate(w, "activated.html", dat)
	}
}

func (app *application) welcomeHandler(w http.ResponseWriter, r *http.Request) {

	// Copy the data from the request body into a new User struct. Notice also that we
	// set the Activated field to false, which isn't strictly necessary because the
	// Activated field will have the zero-value of false by default. But setting this
	// explicitly helps to make our intentions clear to anyone reading the code.
	user := &data.User{
		Name:      r.FormValue("txt"),
		Email:     r.FormValue("email"),
		Activated: false,
	}

	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err := user.Password.Set(r.FormValue("pswd"))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Insert the user data into the database.
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// After the user record has been created in the database, generate a new activation
	// token for the user.
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.background(func() {
		// As there are now multiple pieces of data that we want to pass to our email
		// templates, we create a map to act as a 'holding structure' for the data. This
		// contains the plaintext version of the activation token for the user, along
		// with their ID.
		/* data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		} */
		data := "http://localhost:8000/activated/"
		data += token.Plaintext
		fmt.Print(data)
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.Print("ERROR", err)
		}
	})

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	} else {
		tpl, err := template.ParseFiles("templates/welcome.html")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		dat := r.FormValue("txt")

		tpl.ExecuteTemplate(w, "welcome.html", dat)
	}

}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "token",
		Value:  "",
		Expires: time.Now().AddDate(0, 0, -1)
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

package main

import (
	"net/http"
)

/* func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("pswd")

	v := validator.New()
	data.ValidateEmail(v, email)
	data.ValidatePasswordPlaintext(v, password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

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

	match, err := user.Password.Matches(password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	expiration := token.Expiry.UTC()
	fmt.Println(expiration)
	// set the cookie to expire in 24 hours
	cookie := http.Cookie{
		Name:     "token",
		Value:    base64.StdEncoding.EncodeToString(token.Hash), // replace this with the actual token value
		Path:     "/",
		Expires:  expiration,
		HttpOnly: true, // set HttpOnly to true to prevent client-side access to the cookie
	}

	http.SetCookie(w, &cookie)
	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	http.Redirect(w, r, "/myacc/tasks", http.StatusSeeOther)

} */

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("authenticate"))

	app.writeJSON(w, http.StatusCreated, envelope{"msg": "data have been created"}, nil)

}

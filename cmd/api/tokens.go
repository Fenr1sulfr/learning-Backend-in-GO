package main

import (
	"errors"
	"net/http"
	"time"

	"sulfur.test.net/internal/data"
	"sulfur.test.net/internal/data/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePassowrdPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user, err := app.models.Users.GetByEmail(input.Email)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecordFound):
			app.invalidCredentialResponse(w, r)
		default:
			app.serverErrorRespone(w, r, err)
		}
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorRespone(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialResponse(w, r)
		return
	}
	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorRespone(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorRespone(w, r, err)
		return
	}

}

package main

import (
	"fmt"
	"net/http"
)

func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	app.errorRespone(w, r, http.StatusForbidden, message)
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	app.errorRespone(w, r, http.StatusUnauthorized, message)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorRespone(w, r, http.StatusConflict, message)
}
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorRespone(w, r, http.StatusUnprocessableEntity, errors)
}
func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorRespone(w, r, http.StatusUnauthorized, message)
}
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	app.errorRespone(w, r, http.StatusForbidden, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorRespone(w, r, http.StatusBadRequest, err.Error())
}
func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}
func (app *application) invalidCredentialResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorRespone(w, r, http.StatusUnauthorized, message)
}
func (app *application) errorRespone(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)

	}
}

func (app *application) serverErrorRespone(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorRespone(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorRespone(w, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorRespone(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) rateLimitExceedResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorRespone(w, r, http.StatusTooManyRequests, message)
}

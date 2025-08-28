package main

import (
	"fmt"
	service "goidc_gd/internal"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const grantType = "authorization_code"
const redirectURL = "http://localhost:4000/callback"

const authURL = `https://accounts.google.com/o/oauth2/v2/auth?
  client_id=%s
  &redirect_uri=%s
  &response_type=code
  &scope=https://www.googleapis.com/auth/drive.file
  &state=%s
  &nonce=%s`

func (app *application) displayLogin(w http.ResponseWriter, r *http.Request) {
	w.Write(loginPage)
}

func (app *application) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Generate a cryptographically secure random string or UUID as the state
	// Store this state securely and tied to the user’s browser session(secured http only cookie)
	// redirect the user to idp login set state to the unique id just generated

	state := uuid.New().String()
	nonce := uuid.New().String()
	app.sessionManager.Put(r.Context(), "oidc_state", state)
	app.sessionManager.Put(r.Context(), "nonce", nonce)
	http.Redirect(w, r, fmt.Sprintf(authURL, app.clientID, redirectURL, state, nonce), http.StatusSeeOther)
}

func (app *application) handleLogout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) handleCallback(w http.ResponseWriter, r *http.Request) {
	// read the state from the query string
	// read the uuid from the cookie
	// verify they both match if not reject the request
	// if they match then exchange the code for the tokens

	// decode the id_token to read the user info
	// verify if user exists in the database if not add them to database
	// also store the access token in a database using a session manager like alexedwards/scs for subsequent interactions
	// add a secured cookie to the user's browser indicating they are logged in
	// check if the user's google drive account have a folder called "goidc_gd"
	// if not create it
	// redirect the user to the homepage

	stateFromURLParam := chi.URLParam(r, "state")
	stateFromSession, ok := app.sessionManager.Get(r.Context(), "oidc_state").(string)
	if !ok || stateFromSession != stateFromURLParam {
		app.errorResponse(w, r, http.StatusForbidden, "Invalid State")
		return
	}
	code := chi.URLParam(r, "code")

	// exchange code for the access tokens
	tokens, err := service.GetAccessTokens(app.clientID, app.clientSecret, code, redirectURL, grantType)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) homepage(w http.ResponseWriter, r *http.Request) {
	// verify if the user is logged in
	// if not redirect to the login page
	// if yes fetch all the files in the app default folder using the user's access tokens
	// and render the template and send the response
}

func (app *application) handleFileUpload(w http.ResponseWriter, r *http.Request) {
	// verify if the user is logged in
	// if not redirect to the login page
	// fetch the access tokens and handle multipart file upload to the app default folder in the user's google drive using the access tokens
}

func (app *application) handleFileDownload(w http.ResponseWriter, r *http.Request) {
	// verify if the user is logged in
	// if not redirect to the login page
	// fetch the access tokens and download file from the user's google drive using the access tokens
}

func (app *application) handleFileDelete(w http.ResponseWriter, r *http.Request) {
	// verify if the user is logged in
	// if not redirect to the login page
	// fetch the access tokens and delete file from the user's google drive using the access tokens
}

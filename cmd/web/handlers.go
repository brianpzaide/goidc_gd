package main

import (
	"net/http"
)

func (app *application) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Generate a cryptographically secure random string or UUID as the state
	// Store this state securely and tied to the user’s browser session(secured http only cookie)
	// redirect the user to idp login set state to the unique id just generated
}

func (app *application) handleLogout(w http.ResponseWriter, r *http.Request) {
	// invalidate the session
}

func (app *application) handleCallback(w http.ResponseWriter, r *http.Request) {
	// read the state from the query string
	// read the uuid from the cookie
	// verify they both match if not reject the request
	// if they match exchange the code for the tokens
	// decode the id_token to read the user info
	// verify if user exists in the database if not add them to database
	// also store the access token in a database using a session manager like alexedwards/scs for subsequent interactions
	// add a secured cookie to the user's browser indicating they are logged in
	// check if the user's google drive account have a folder called "goidc_gd"
	// if not create it
	// redirect the user to the homepage
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

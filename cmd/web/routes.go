package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func routes(app *application) http.Handler {
	router := chi.NewRouter()

	router.NotFound(app.notFoundResponse)
	router.MethodNotAllowed(app.methodNotAllowedResponse)

	router.Get("/login", app.handleLogin)
	router.Get("/logout", app.handleLogout)
	router.Get("/callback", app.handleCallback)
	router.Get("/", app.homepage)
	router.Route("/files", func(r chi.Router) {
		r.Post("/", app.handleFileUpload)
		r.Delete("/{file_id}", app.handleFileDelete)
		r.Get("/{file_id}", app.handleFileDownload)
	})

	return router
}

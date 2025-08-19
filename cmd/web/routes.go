package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func routes(app *application) http.Handler {
	router := chi.NewRouter()

	router.NotFound(app.notFoundResponse)
	router.MethodNotAllowed(app.methodNotAllowedResponse)

	router.Get("/login", app.displayLogin)
	router.Get("/auth/{provider}", app.handleLogin)
	router.Get("/callback", app.handleCallback)

	// protected endpoints
	router.Route("", func(r1 chi.Router) {
		r1.Use(app.authenticate)

		r1.Get("/", app.homepage)
		r1.Post("/logout", app.handleLogout)
		r1.Route("/files", func(r2 chi.Router) {
			r2.Post("/", app.handleFileUpload)
			r2.Delete("/{file_id}", app.handleFileDelete)
			r2.Get("/{file_id}", app.handleFileDownload)
		})
	})

	return router
}

package quotes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Start() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/quotes/", func(r chi.Router) {
		r.Route("/session/", func(r chi.Router) {
			// TODO: create session
			// TODO: join session
		})
		r.Route("/collection/", func(r chi.Router) {
			r.Post("/create/", createCollection)
			// TODO: create collection
			// TODO: modify collection
			// TODO: delete collection

		})
		r.Route("/quote/", func(r chi.Router) {
			// TODO: add quote to collection
			r.Post("/add/", addQuote)
			// TODO: modify quote
			// TODO: delete quote
		})
	})

	http.ListenAndServe(":3000", r)
}

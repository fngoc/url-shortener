package server

import (
	"github.com/fngoc/url-shortener/cmd/shortener/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

const port string = ":8080"

// Run функция будет полезна при инициализации зависимостей сервера перед запуском
func Run() error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.PostWebhook)
		r.Get("/{id}", handlers.GetWebhook)
	})

	return http.ListenAndServe(port, r)
}

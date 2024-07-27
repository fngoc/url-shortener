package server

import (
	"github.com/fngoc/url-shortener/cmd/shortener/config"
	"github.com/fngoc/url-shortener/cmd/shortener/handlers"
	"github.com/fngoc/url-shortener/internal/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Run функция будет полезна при инициализации зависимостей сервера перед запуском
func Run() error {
	if err := logger.Initialize(); err != nil {
		return err
	}

	logger.Log.Info("Starting server")

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", logger.RequestLogger(handlers.PostSaveWebhook))
		r.Get("/{id}", logger.RequestLogger(handlers.GetRedirectWebhook))
		r.Route("/api", func(r chi.Router) {
			r.Post("/shorten", logger.RequestLogger(handlers.PostShortenWebhook))
		})
	})

	return http.ListenAndServe(config.Flags.ServerAddress, r)
}

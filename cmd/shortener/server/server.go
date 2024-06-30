package server

import (
	"github.com/fngoc/url-shortener/cmd/shortener/handlers"
	"net/http"
)

const port string = ":8080"

// Run функция будет полезна при инициализации зависимостей сервера перед запуском
func Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.PostWebhook)
	mux.HandleFunc("/{id}", handlers.GetWebhook)

	return http.ListenAndServe(port, mux)
}

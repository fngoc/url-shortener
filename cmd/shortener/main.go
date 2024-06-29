package main

import (
	"github.com/fngoc/url-shortener/internal/app"
	"io"
	"net/http"
	"strings"
)

const host string = "http://localhost"
const port string = ":8080"

type LocalStore map[string]string

var store = make(LocalStore)

// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", postWebhook)
	mux.HandleFunc("/{id}", getWebhook)

	return http.ListenAndServe(port, mux)
}

// функция getWebhook — обработчик GET HTTP-запроса
func getWebhook(w http.ResponseWriter, r *http.Request) {
	// разрешаем только GET-запросы
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// читаем url
	id := strings.TrimPrefix(r.URL.Path, "/")
	// проверяем id
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	// получаем url из локальной памяти
	url, ok := store[id]
	// проверяем наличие url в локальной памяти
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// редиректим на url
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// функция postWebhook — обработчик POST HTTP-запроса
func postWebhook(w http.ResponseWriter, r *http.Request) {
	// разрешаем только POST-запросы с Content-Type: text/plain
	if r.Method != http.MethodPost || !strings.Contains(r.Header.Get("Content-Type"), "text/plain") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// читаем все body
	b, _ := io.ReadAll(r.Body)

	if len(b) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// генерируем строку
	id := app.GenerateString(8)
	// сохраняем в локальную память
	store[id] = string(b)

	// установим правильный заголовок для типа данных
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	// пока установим ответ-заглушку, без проверки ошибок
	_, _ = w.Write([]byte(host + port + "/" + id))
}

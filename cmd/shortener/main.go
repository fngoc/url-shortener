package main

import (
	"github.com/fngoc/url-shortener/internal/app"
	"io/ioutil"
	"net/http"
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
	// разрешаем только GET-запросы с Content-Type: text/plain
	if r.Method != "GET" || r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// читаем url
	b := []byte(r.URL.String())
	// получаем url из локальной памяти
	url, ok := store[string(b[1:])]
	// проверяем наличие url в локальной памяти
	if !ok {
		w.WriteHeader(http.StatusNotFound)
	}
	// редиректим на url
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// функция postWebhook — обработчик POST HTTP-запроса
func postWebhook(w http.ResponseWriter, r *http.Request) {
	// разрешаем только POST-запросы с Content-Type: text/plain
	if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// читаем все body
	b, _ := ioutil.ReadAll(r.Body)
	// генерируем строку
	str := app.GenerateString(10)
	// сохраняем в локальную память
	store[str] = string(b)

	// установим правильный заголовок для типа данных
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	// пока установим ответ-заглушку, без проверки ошибок
	_, _ = w.Write([]byte(host + port + "/" + str))
}

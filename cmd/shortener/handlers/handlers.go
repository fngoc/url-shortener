package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/fngoc/url-shortener/cmd/shortener/config"
	"github.com/fngoc/url-shortener/cmd/shortener/storage"
	"github.com/fngoc/url-shortener/internal/models"
	"github.com/fngoc/url-shortener/internal/utils"
	"io"
	"net/http"
	"strings"
)

// GetRedirectWebhook функция обработчик GET HTTP-запроса
func GetRedirectWebhook(w http.ResponseWriter, r *http.Request) {
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
	url, err := storage.Store.GetData(id)
	// проверяем наличие url в локальной памяти
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// редиректим на url
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// PostSaveWebhook функция обработчик POST HTTP-запроса
func PostSaveWebhook(w http.ResponseWriter, r *http.Request) {
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
	id := utils.GenerateString(8)
	// сохраняем в локальную память
	err := storage.Store.SaveData(id, string(b))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// установим правильный заголовок для типа данных
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	// пока установим ответ-заглушку, без проверки ошибок
	_, _ = w.Write([]byte(config.Flags.BaseResultAddress + "/" + id))
}

// PostShortenWebhook функция обработчик POST HTTP-запроса
func PostShortenWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// читаем все body
	dec := json.NewDecoder(r.Body)
	var req models.Request
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// генерируем строку
	id := utils.GenerateString(8)
	// сохраняем в локальную память
	err := storage.Store.SaveData(id, req.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	buf := bytes.Buffer{}
	encode := json.NewEncoder(&buf)
	if err := encode.Encode(models.Response{Result: config.Flags.BaseResultAddress + "/" + id}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// установим правильный заголовок для типа данных
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write(buf.Bytes())
}

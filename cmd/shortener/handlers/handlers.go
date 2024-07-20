package handlers

import (
	"github.com/fngoc/url-shortener/cmd/shortener/config"
	"github.com/fngoc/url-shortener/cmd/shortener/storage"
	"github.com/fngoc/url-shortener/internal/utils"
	"io"
	"net/http"
	"strings"
)

// GetWebhook функция обработчик GET HTTP-запроса
func GetWebhook(w http.ResponseWriter, r *http.Request) {
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

// PostWebhook функция обработчик POST HTTP-запроса
func PostWebhook(w http.ResponseWriter, r *http.Request) {
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

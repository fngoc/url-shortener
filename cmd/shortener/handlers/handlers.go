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
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	url, err := storage.Store.GetData(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// PostSaveWebhook функция обработчик POST HTTP-запроса
func PostSaveWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || (!strings.Contains(r.Header.Get("Content-Type"), "text/plain") &&
		!strings.Contains(r.Header.Get("Content-Type"), "application/x-gzip")) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, _ := io.ReadAll(r.Body)

	if len(b) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := utils.GenerateString(8)
	err := storage.Store.SaveData(id, string(b))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(config.Flags.BaseResultAddress + "/" + id))
}

// PostShortenWebhook функция обработчик POST HTTP-запроса
func PostShortenWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(r.Body)
	var req models.Request
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := utils.GenerateString(8)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write(buf.Bytes())
}

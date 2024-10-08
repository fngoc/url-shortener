package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/fngoc/url-shortener/cmd/shortener/config"
	"github.com/fngoc/url-shortener/cmd/shortener/constants"
	"github.com/fngoc/url-shortener/cmd/shortener/storage"
	"github.com/fngoc/url-shortener/internal/logger"
	"github.com/fngoc/url-shortener/internal/models"
	"github.com/fngoc/url-shortener/internal/utils"
	"github.com/jackc/pgerrcode"
	"io"
	"net/http"
	"strings"
)

type deleteJob struct {
	userID int
	url    string
}

// GetRedirectWebhook функция обработчик GET HTTP-запроса
func GetRedirectWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url, err := storage.Store.GetData(r.Context(), id)
	if err != nil {
		var deleteErr *storage.DBDeleteError
		if errors.As(err, &deleteErr) {
			w.WriteHeader(http.StatusGone)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GetUrlsWebhook функция обработчик GET HTTP-запроса для получения всех urls
func GetUrlsWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie, _ := r.Cookie(CookieName)
	if cookie == nil {
		authHeader := r.Header.Get("Authorization")
		if _, err := GetUserID(authHeader); err != nil {
			logger.Log.Warn("Token in 'Authorization' header is not valid")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	userID := r.Context().Value(constants.UserIDKey).(int)
	authCtx := context.WithValue(r.Context(), constants.UserIDKey, userID)

	urls, err := storage.Store.GetAllData(authCtx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	buf := bytes.Buffer{}
	encode := json.NewEncoder(&buf)
	if err := encode.Encode(urls); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

// DeleteUrlsWebhook функция обработчик DELETE HTTP-запроса для удаления urls
func DeleteUrlsWebhook(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(constants.UserIDKey).(int)
	var IDs []string

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&IDs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	numWorkers := 3
	numJobs := len(IDs)

	jobs := make(chan deleteJob, numJobs)
	defer close(jobs)

	for w := 0; w < numWorkers; w++ {
		go deleteWorker(jobs)
	}

	for j := 0; j < numJobs; j++ {
		item := deleteJob{
			userID: userID,
			url:    IDs[j],
		}
		jobs <- item
	}

	w.WriteHeader(http.StatusAccepted)
}

func deleteWorker(jobs <-chan deleteJob) {
	for j := range jobs {
		storage.Store.DeleteData(j.userID, j.url)
	}
}

// PostSaveWebhook функция обработчик POST HTTP-запроса
func PostSaveWebhook(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	allowedTextPlan := strings.Contains(contentType, "text/plain")
	gzipTextPlan := strings.Contains(contentType, "gzip")

	if r.Method != http.MethodPost || (!allowedTextPlan && !gzipTextPlan) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, _ := io.ReadAll(r.Body)

	if len(b) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := utils.GenerateString(8)
	err := storage.Store.SaveData(r.Context(), id, string(b))
	if err != nil {
		var dbErr *storage.DBError
		if errors.As(err, &dbErr) && pgerrcode.IsIntegrityConstraintViolation(dbErr.Err.Code) {
			setResponsePostSaveWebhook(w, http.StatusConflict, dbErr.ShortURL)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	setResponsePostSaveWebhook(w, http.StatusCreated, id)
}

// setResponsePostSaveWebhook устанавливает ответ для PostSaveWebhook
func setResponsePostSaveWebhook(w http.ResponseWriter, statusCode int, id string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(config.Flags.BaseResultAddress + "/" + id))
}

// PostShortenWebhook функция обработчик POST HTTP-запроса
func PostShortenWebhook(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	allowedApplicationJSON := strings.Contains(contentType, "application/json")
	gzipTextPlan := strings.Contains(contentType, "gzip")

	if r.Method != http.MethodPost || (!allowedApplicationJSON && !gzipTextPlan) {
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
	err := storage.Store.SaveData(r.Context(), id, req.URL)
	if err != nil {
		var dbErr *storage.DBError
		if errors.As(err, &dbErr) && pgerrcode.IsIntegrityConstraintViolation(dbErr.Err.Code) {
			id = dbErr.ShortURL

			buf := bytes.Buffer{}
			encode := json.NewEncoder(&buf)
			if err := encode.Encode(models.Response{Result: config.Flags.BaseResultAddress + "/" + id}); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			setResponsePostShortenWebhook(w, http.StatusConflict, buf)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	buf := bytes.Buffer{}
	encode := json.NewEncoder(&buf)
	if err := encode.Encode(models.Response{Result: config.Flags.BaseResultAddress + "/" + id}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	setResponsePostShortenWebhook(w, http.StatusCreated, buf)
}

// setResponsePostShortenWebhook устанавливает ответ для PostShortenWebhook
func setResponsePostShortenWebhook(w http.ResponseWriter, statusCode int, buf bytes.Buffer) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(buf.Bytes())
}

// PostShortenBatchWebhook функция обработчик POST HTTP-запроса для сохранения данных бачами
func PostShortenBatchWebhook(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	allowedApplicationJSON := strings.Contains(contentType, "application/json")
	gzipTextPlan := strings.Contains(contentType, "gzip")

	if r.Method != http.MethodPost || (!allowedApplicationJSON && !gzipTextPlan) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dec := json.NewDecoder(r.Body)
	var req []models.RequestBatch
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var resp = make([]models.ResponseBatch, 0, len(req))
	for _, v := range req {
		id := utils.GenerateString(8)
		err := storage.Store.SaveData(r.Context(), id, v.OriginalURL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		resp = append(resp, models.ResponseBatch{
			CorrelationID: v.CorrelationID,
			ShortURL:      config.Flags.BaseResultAddress + "/" + id,
		})
	}

	buf := bytes.Buffer{}
	encode := json.NewEncoder(&buf)
	if err := encode.Encode(resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write(buf.Bytes())
}

// CheckConnection функция обработчик GET HTTP-запроса для проверки соединения с БД
func CheckConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !storage.CustomPing() {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

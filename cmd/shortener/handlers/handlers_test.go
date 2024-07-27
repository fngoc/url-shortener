package handlers

import (
	"encoding/json"
	"github.com/fngoc/url-shortener/cmd/shortener/storage"
	"github.com/fngoc/url-shortener/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockLocalStore storage.LocalStore

func TestGetRedirectWebhook(t *testing.T) {
	type want struct {
		statusCode int
		expectBody bool
		body       string
	}
	tests := []struct {
		name       string
		method     string
		requestURL string
		store      MockLocalStore
		want       want
	}{
		{
			"200 code test",
			"GET",
			"/testKeys",
			MockLocalStore{
				"testKeys": "https://google.com",
			},
			want{
				statusCode: http.StatusTemporaryRedirect,
				expectBody: true,
				body:       "https://google.com",
			},
		},
		{
			"not Get method code test",
			"POST",
			"/testKeys",
			nil,
			want{
				statusCode: http.StatusBadRequest,
				expectBody: false,
				body:       "",
			},
		},
		{
			"empty id in url test",
			"GET",
			"/",
			nil,
			want{
				statusCode: http.StatusBadRequest,
				expectBody: false,
				body:       "",
			},
		},
		{
			"empty id in store test",
			"GET",
			"/",
			MockLocalStore{},
			want{
				statusCode: http.StatusBadRequest,
				expectBody: false,
				body:       "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.Store = storage.LocalStore(tt.store)

			request := httptest.NewRequest(tt.method, tt.requestURL, nil)
			w := httptest.NewRecorder()

			GetRedirectWebhook(w, request)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.want.statusCode, res.StatusCode)

			if tt.want.expectBody {
				assert.NotEmpty(t, res.Body)
			}
		})
	}
}

func TestPostSaveWebhook(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		expectBody  bool
	}
	tests := []struct {
		name        string
		method      string
		body        string
		contentType string
		want        want
	}{
		{
			"200 code test",
			"POST",
			"https://ya.com",
			"text/plain",
			want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				expectBody:  true,
			},
		},
		{
			"not POST method test",
			"GET",
			"asdasd",
			"text/plain",
			want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
				expectBody:  false,
			},
		},
		{
			"not text/plain test",
			"POST",
			"https://google.com",
			"application/json",
			want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
				expectBody:  false,
			},
		},
		{
			"empty body test",
			"POST",
			"",
			"text/plain",
			want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
				expectBody:  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			request.Header.Add("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			PostSaveWebhook(w, request)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.want.statusCode, res.StatusCode)
			require.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			if tt.want.expectBody {
				assert.NotEmpty(t, res.Body)
			}
		})
	}
}

func TestPostShortenWebhook(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		expectBody  bool
	}
	tests := []struct {
		name        string
		method      string
		body        string
		contentType string
		want        want
	}{
		{
			"201 code test",
			"POST",
			"{\"url\":\"https://ya.ru\"}",
			"application/json",
			want{
				contentType: "application/json",
				statusCode:  http.StatusCreated,
				expectBody:  true,
			},
		},
		{
			"not POST method test",
			"GET",
			"asdasd",
			"application/json",
			want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
				expectBody:  false,
			},
		},
		{
			"not application/json test",
			"POST",
			"https://google.com",
			"text/plan",
			want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
				expectBody:  false,
			},
		},
		{
			"empty body test",
			"POST",
			"",
			"text/plain",
			want{
				contentType: "",
				statusCode:  http.StatusBadRequest,
				expectBody:  false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/api/shorten", strings.NewReader(tt.body))
			request.Header.Add("Content-Type", tt.contentType)
			w := httptest.NewRecorder()

			PostShortenWebhook(w, request)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.want.statusCode, res.StatusCode)
			require.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			if tt.want.expectBody {
				assert.NotEmpty(t, res.Body)

				dec := json.NewDecoder(res.Body)
				var resp models.Response
				if err := dec.Decode(&resp); err != nil {
					return
				}
				assert.NotEmpty(t, resp.Result)
			}
		})
	}
}

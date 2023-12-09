package handlers

import (
	"bytes"
	"context"
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestURLShortenerHandler_PostHandler1(t *testing.T) {
	type fields struct {
		URLRepository storage.URLRepository
		target        string
		url           string
	}
	type want struct {
		code        int
		body        string
		contentType string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "positive test",
			fields: fields{
				URLRepository: storage.NewURLStorage(),
				target:        "/",
				url:           "https://practicum.yandex.ru/",
			},
			want: want{
				code:        http.StatusCreated,
				body:        "http://example.com/CKj87ajs",
				contentType: "text/plain",
			},
		},
		{
			name: "negative test",
			fields: fields{
				URLRepository: storage.NewURLStorage(),
				target:        "/",
				url:           "",
			},
			want: want{
				code:        http.StatusBadRequest,
				body:        "",
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &URLShortenerHandler{
				URLRepository: test.fields.URLRepository,
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, test.fields.target, bytes.NewReader([]byte(test.fields.url)))
			handler.PostHandler(w, r)

			res := w.Result()
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			str := string(resBody)

			require.NoError(t, err)
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.body, str)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestURLShortenerHandler_GetHandler(t *testing.T) {
	type fields struct {
		URLRepository storage.URLRepository
		target        string
		url           string
		store         bool
	}
	type want struct {
		code        int
		contentType string
		location    string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "positive test",
			fields: fields{
				URLRepository: storage.NewURLStorage(),
				target:        "CKj87ajs",
				url:           "https://practicum.yandex.ru/",
				store:         true,
			},
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/plain",
				location:    "https://practicum.yandex.ru/",
			},
		},
		{
			name: "negwtive test",
			fields: fields{
				URLRepository: storage.NewURLStorage(),
				target:        "CKj87ajs",
				url:           "https://practicum.yandex.ru/",
				store:         false,
			},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
				location:    "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := &URLShortenerHandler{
				URLRepository: test.fields.URLRepository,
			}

			if test.fields.store {
				handler.URLRepository.AddURL(test.fields.target, test.fields.url)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/{id}", http.NoBody)

			routeContext := chi.NewRouteContext()
			routeContext.URLParams.Add("id", test.fields.target)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routeContext))
			handler.GetHandler(w, r)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.location, res.Header.Get("Location"))
		})
	}
}

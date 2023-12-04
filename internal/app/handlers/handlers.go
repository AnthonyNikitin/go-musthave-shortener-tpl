package handlers

import (
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/hasher"
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type URLShortenerHandler struct {
	URLStorage *storage.URLStorage
}

func NewURLShortenerHandler() *URLShortenerHandler {
	return &URLShortenerHandler{
		URLStorage: storage.NewURLStorage(),
	}
}

func (handler *URLShortenerHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := string(body)
	shortLink := hasher.GetShortLink(res)
	handler.URLStorage.Urls[shortLink] = res

	w.WriteHeader(http.StatusCreated)
	host := r.Host
	_, err = w.Write([]byte("http://" + host + "/" + shortLink))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (handler *URLShortenerHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	url, ok := handler.URLStorage.Urls[id]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

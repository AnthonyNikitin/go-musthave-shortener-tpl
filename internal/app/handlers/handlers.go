package handlers

import (
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/hasher"
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type URLShortenerHandler struct {
	UrlRepository storage.URLRepository
}

func NewURLShortenerHandler() *URLShortenerHandler {
	return &URLShortenerHandler{
		UrlRepository: storage.NewURLStorage(),
	}
}

func (handler *URLShortenerHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	w.Header().Set("Content-Type", "text/plain")

	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := string(body)
	shortLink, err := hasher.GetShortLink(res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	handler.UrlRepository.AddURL(shortLink, res)

	w.WriteHeader(http.StatusCreated)
	host := r.Host

	_, err = w.Write([]byte("http://" + host + "/" + shortLink))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (handler *URLShortenerHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	url, ok := handler.UrlRepository.GetURL(id)

	w.Header().Set("Content-Type", "text/plain")

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

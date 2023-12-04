package main

import (
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/hasher"
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
)

var urlStorage = storage.NewUrlStorage()

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/", postHandler)
	r.Get("/{id}", getHandler)

	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := string(body)
	shortLink := hasher.GetShortLink(res)
	urlStorage.Urls[shortLink] = res

	w.WriteHeader(http.StatusCreated)
	host := r.Host
	_, err = w.Write([]byte("http://" + host + "/" + shortLink))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	url, ok := urlStorage.Urls[id]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

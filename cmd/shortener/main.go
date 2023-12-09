package main

import (
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("text/plain"))

	urlShortenerHandler := handlers.NewURLShortenerHandler()

	r.Post("/", urlShortenerHandler.PostHandler)
	r.Get("/{id}", urlShortenerHandler.GetHandler)

	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}
}

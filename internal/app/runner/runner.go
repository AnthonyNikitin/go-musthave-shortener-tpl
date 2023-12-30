package runner

import (
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/config"
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/middlewares"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RunApplication() error {

	c := config.NewConfiguration()
	c.ParseConfiguration()

	r := chi.NewRouter()
	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.GzipMiddleware)

	urlShortenerHandler, err := handlers.NewURLShortenerHandler(c.BaseResponseURL, c.FileStoragePath)

	if err != nil {
		panic(err)
	}

	r.Post("/", urlShortenerHandler.PostHandler)
	r.Get("/{id}", urlShortenerHandler.GetHandler)

	r.Route("/api", func(r chi.Router) {
		r.Post("/shorten", urlShortenerHandler.PostShortenHandler)
	})

	err = http.ListenAndServe(c.Address, r)
	if err != nil {
		return err
	}

	return nil
}

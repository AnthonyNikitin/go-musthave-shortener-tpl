package runner

import (
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/config"
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func RunApplication() error {
	c := config.NewConfiguration()
	c.ParseConfiguration()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("text/plain"))

	urlShortenerHandler := handlers.NewURLShortenerHandler(c.BaseResponseURL)

	r.Post("/", urlShortenerHandler.PostHandler)
	r.Get("/{id}", urlShortenerHandler.GetHandler)

	err := http.ListenAndServe(c.Address, r)
	if err != nil {
		return err
	}

	return nil
}

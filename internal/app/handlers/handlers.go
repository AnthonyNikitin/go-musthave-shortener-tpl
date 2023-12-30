package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/hasher"
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type URLShortenerHandler struct {
	URLRepository   storage.URLRepository
	BaseResponseURL string
}

func NewURLShortenerHandler(baseResponseURL, fileStoragePath string) (*URLShortenerHandler, error) {

	urlRepository, err := storage.NewURLStorage(fileStoragePath)
	if err != nil {
		return nil, err
	}

	return &URLShortenerHandler{
		URLRepository:   urlRepository,
		BaseResponseURL: baseResponseURL,
	}, nil
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
		fmt.Println(err.Error())
		return
	}

	err = handler.URLRepository.AddURL(shortLink, res)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(handler.BaseResponseURL + shortLink))
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

func (handler *URLShortenerHandler) PostShortenHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var request ShortenRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	if len(request.URL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("empty request url")
		return
	}

	shortLink, err := hasher.GetShortLink(request.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	err = handler.URLRepository.AddURL(shortLink, request.URL)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}

	result := handler.BaseResponseURL + shortLink
	response := ShortenResponse{
		Result: result,
	}

	output, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(output)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (handler *URLShortenerHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	url, ok := handler.URLRepository.GetURL(id)

	w.Header().Set("Content-Type", "text/plain")

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

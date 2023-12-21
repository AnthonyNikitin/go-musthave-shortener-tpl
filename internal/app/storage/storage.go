package storage

import "sync"

type URLRepository interface {
	AddURL(shortLink, originalLink string)
	GetURL(shortLink string) (string, bool)
}

type URLStorage struct {
	mu   sync.Mutex
	Urls map[string]string
}

func NewURLStorage() *URLStorage {
	return &URLStorage{
		Urls: make(map[string]string),
	}
}

func (s *URLStorage) AddURL(shortLink, originalLink string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Urls[shortLink] = originalLink
}

func (s *URLStorage) GetURL(shortLink string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	url, ok := s.Urls[shortLink]

	return url, ok
}

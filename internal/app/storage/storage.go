package storage

type URLRepository interface {
	AddURL(shortLink, originalLink string)
	GetURL(shortLink string) (string, bool)
}

type URLStorage struct {
	Urls map[string]string
}

func NewURLStorage() *URLStorage {
	return &URLStorage{
		Urls: make(map[string]string),
	}
}

func (s *URLStorage) AddURL(shortLink, originalLink string) {
	s.Urls[shortLink] = originalLink
}

func (s *URLStorage) GetURL(shortLink string) (string, bool) {
	url, ok := s.Urls[shortLink]

	return url, ok
}

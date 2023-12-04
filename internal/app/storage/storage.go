package storage

type URLStorage struct {
	Urls map[string]string
}

func NewURLStorage() *URLStorage {
	return &URLStorage{
		Urls: make(map[string]string),
	}
}

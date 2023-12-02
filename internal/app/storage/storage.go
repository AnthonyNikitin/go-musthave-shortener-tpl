package storage

type UrlStorage struct {
	Urls map[string]string
}

func NewUrlStorage() *UrlStorage {
	return &UrlStorage{
		Urls: make(map[string]string),
	}
}

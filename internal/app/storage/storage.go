package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type URLRepository interface {
	AddURL(shortLink, originalLink string) error
	GetURL(shortLink string) (string, bool)
}

type URLStorage struct {
	mu              sync.Mutex
	Urls            map[string]string
	fileStoragePath string
}

type URLEntity struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewURLStorage(fileStoragePath string) (*URLStorage, error) {

	urls := make(map[string]string)
	if len(fileStoragePath) != 0 {
		//fileStoragePath = strings.TrimPrefix(fileStoragePath, "/")

		fileStorage, err := createFileStorage(fileStoragePath)
		defer func(fileStorage *os.File) {
			err = fileStorage.Close()
		}(fileStorage)

		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(fileStorage)
		var urlEntity URLEntity
		for scanner.Scan() {

			bytes := scanner.Bytes()
			if len(bytes) == 0 {
				continue
			}
			err := json.Unmarshal(bytes, &urlEntity)

			if err != nil {
				return nil, err
			}

			urls[urlEntity.ShortURL] = urlEntity.OriginalURL
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	return &URLStorage{
		Urls:            urls,
		fileStoragePath: fileStoragePath,
	}, nil
}

func (s *URLStorage) AddURL(shortLink, originalLink string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.fileStoragePath) > 0 {
		err := s.writeToFile(shortLink, originalLink)
		if err != nil {
			return err
		}
	}

	s.Urls[shortLink] = originalLink

	return nil
}

func (s *URLStorage) GetURL(shortLink string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	url, ok := s.Urls[shortLink]

	return url, ok
}

func createFileStorage(p string) (*os.File, error) {
	if _, err := os.Stat(p); err == nil {
		file, err := os.Open(p)
		if err != nil {
			return nil, err
		}

		return file, nil
	}

	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func (s *URLStorage) writeToFile(shortLink, originalLink string) error {

	file, err := os.OpenFile(s.fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	buffer := bufio.NewWriter(file)

	urlEntity := URLEntity{
		ShortURL:    shortLink,
		OriginalURL: originalLink,
	}

	bytes, err := json.Marshal(urlEntity)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(string(bytes) + "\n")
	if err != nil {
		return err
	}

	if err := buffer.Flush(); err != nil {
		return err
	}

	return nil
}

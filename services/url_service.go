package services

import (
	"errors"
	"time"

	"github.com/Shashankm886/url-shortener/storage"
	"github.com/Shashankm886/url-shortener/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLService struct {
	store *storage.MongoDBStore
}

func NewURLService(store *storage.MongoDBStore) *URLService {
	return &URLService{store: store}
}

func (s *URLService) GenerateUniqueShortURL() string {
	for {
		shortURL := utils.GenerateShortURL()
		_, err := s.store.Get(shortURL)
		if err != nil {
			return shortURL
		}
	}
}

func (s *URLService) ShortenURL(longURL string, expirySeconds int) (string, error) {
	if expirySeconds <= 0 {
		expirySeconds = 3600
	}

	shortURL := s.GenerateUniqueShortURL()
	expiryTime := time.Now().Add(time.Duration(expirySeconds) * time.Second)

	err := s.store.Save(shortURL, longURL, expiryTime)
	if err != nil {
		return "", errors.New("failed to save URL")
	}
	return shortURL, nil
}

func (s *URLService) GetLongURL(shortURL string) (string, bool) {
	data, err := s.store.Get(shortURL)
	if err != nil {
		return "", false
	}

	expirationVal, ok := (*data)["expiration_time"]
	if !ok {
		return "", false
	}

	expirationPrimitive, ok := expirationVal.(primitive.DateTime)
	if !ok {
		return "", false
	}

	expirationTime := expirationPrimitive.Time()

	if time.Now().After(expirationTime) {
		_ = s.store.Delete(shortURL)
		return "", false
	}

	err = s.store.IncrementUsage(shortURL)
	if err != nil {
		return "", false
	}

	longURL, _ := (*data)["long_url"].(string)
	return longURL, true
}

func (s *URLService) GetUsage(shortURL string) (int, error) {
	data, err := s.store.Get(shortURL)
	print("stop 0")
	if err != nil {
		return 0, err
	}
	print("stop 1")
	usage, ok := (*data)["usage"].(int32)
	if !ok {
		return 0, errors.New("invalid usage data")
	}
	print("stop 2")

	return int(usage), nil
}

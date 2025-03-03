package controllers

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/Shashankm886/url-shortener/services"
	"github.com/Shashankm886/url-shortener/storage"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var urlService *services.URLService

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatalf("MONGO_URI is not set in the .env file")
	}

	store, err := storage.NewMongoDBStore(mongoURI, "urlShortenerDB", "urls")
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	urlService = services.NewURLService(store)
}

func isValidURL(input string) bool {
	parsedURL, err := url.ParseRequestURI(input)
	return err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

func ShortenURL(c *gin.Context) {
	var reqBody struct {
		LongURL           string `json:"long_url"`
		ExpirationSeconds int    `json:"expiration_seconds"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if !isValidURL(reqBody.LongURL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid URL"})
		return
	}

	expirySeconds := reqBody.ExpirationSeconds
	if expirySeconds <= 0 {
		expirySeconds = 3600
	}

	shortURL, err := urlService.ShortenURL(reqBody.LongURL, expirySeconds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to shorten URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"long_url": reqBody.LongURL, "short_url": shortURL})
}

func RedirectURL(c *gin.Context) {
	shortURL := c.Param("shortUrl")

	longURL, found := urlService.GetLongURL(shortURL)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, longURL)
}

func GetUsage(c *gin.Context) {
	shortURL := c.Param("shortUrl")

	usage, err := urlService.GetUsage(shortURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"short_url": shortURL, "usage_count": usage})
}

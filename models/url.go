package models

type URL struct {
	LongURL  string `json:"long_url" binding:"required"`
	ShortURL string `json:"short_url"`
}

package domain

import "time"

type Url struct {
	LongUrl   string `json:"longUrl,omitempty"`
	ShortUrl  string `json:"shortUrl,omitempty"`
	CreatedAt time.Time
}

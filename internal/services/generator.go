package services

import (
	"math/rand"
	"time"
)

const shortUrlLength = 10

var alphabet = []byte(`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_`)

func genShort() string {
	rand.Seed(time.Now().UnixNano())

	shortUrl := make([]byte, 10)
	for i := 0; i < shortUrlLength; i++ {
		shortUrl[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(shortUrl)
}

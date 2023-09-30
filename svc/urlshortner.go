package svc

import (
	"math/rand"
	"time"
)

type Mode string

const (
	Random Mode = "Random"
	Phrase Mode = "Phrase"

	charset = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type URLShortner struct {
	mode   Mode
	length int
}

func NewKeywordURLShortner(length int) URLShortner {
	return URLShortner{
		mode:   Phrase,
		length: length,
	}
}

func NewRandomURLShortner(length int) URLShortner {
	return URLShortner{
		mode:   Random,
		length: length,
	}
}

func (u *URLShortner) Shorten(longURL string) string {
	if u.mode == Phrase {
		// TODO: Implement
		return "hello"
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomBytes := make([]byte, u.length)
	for i := range randomBytes {
		randomBytes[i] = charset[r.Intn(len(charset))]
	}
	return string(randomBytes)
}

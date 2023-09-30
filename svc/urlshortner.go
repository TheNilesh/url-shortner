package svc

type Mode string

const (
	Random Mode = "Random"
	Phrase Mode = "Phrase"
)

type URLShortner struct {
	mode    Mode
	length  int
	charset string
}

func NewKeywordURLShortner(length int) URLShortner {
	return URLShortner{
		mode:   Phrase,
		length: length,
	}
}

func NewRandomURLShortner(charset string, length int) URLShortner {
	return URLShortner{
		mode:    Random,
		charset: charset,
		length:  length,
	}
}

func (u *URLShortner) Shorten(longURL string) string {
	// TODO: Actual business logic
	return ""
}

package api

type Mode string

const (
	Random Mode = "RANDOM"
	Phase  Mode = "Phase"
)

type URLShortner struct {
	mode    Mode
	length  int
	charset string
}

func NewKeywordURLShortner(length int) URLShortner {
	return URLShortner{
		mode:   Phase,
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

func (u *URLShortner) Shorten() (string, error) {
	// TODO: Actual business logic
	return "", nil
}

package api

type ShortURL struct {
	ID     string
	URL    string
	Domain string
}

type URLShortner struct {
}

func (u *URLShortner) Shorten() (*ShortURL, error) {
	// TODO: Actual business logic
	return nil, nil
}

func (u *URLShortner) LengthenURL() (string, error) {
	// TODO: Actual business logic
	return "", nil
}

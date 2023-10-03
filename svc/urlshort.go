package svc

import "errors"

var (
	ErrServerError = errors.New("server error")
	ErrConflict    = errors.New("conflict")
)

func CreateShortURL(shortPath string, targetURL string) (string, error) {
	if len(shortPath) > 0 { // isShortPathProvidedInRequest ?
		oldTargetURL, found, err := lookupTargetURL(shortPath)
		if err != nil {
			// return 500
			return "", ErrServerError
		}
		if found {
			if oldTargetURL == targetURL {
				// return OK
				return shortPath, nil
			} else {
				// return 409
				return "", ErrConflict
			}
		}
	}
	existingShortPath, found, err := lookupShortPath(targetURL)
	if err != nil {
		return "", ErrServerError
	}
	if found { // isTargetURLAlreadyShortened
		return existingShortPath, nil
	}
	return doShorten(shortPath, targetURL)
}

func doShorten(shortPath string, targetURL string) (string, error) {
	if len(shortPath) == 0 {
		shortPath = generateShortPath()
	}
	err := saveShortPathToURLMapping(shortPath, targetURL)
	if err != nil {
		return "", ErrServerError
	}
	return shortPath, nil
}

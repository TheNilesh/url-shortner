package svc

import (
	"math/rand"
	"strings"
	"time"

	"github.com/thenilesh/url-shortner/store"
)

type Mode string

const (
	Random Mode = "Random"
	Phrase Mode = "Phrase"

	charset = "abcdefghijklmnopqrstuvwxyz-0123456789"
)

type URLShortner struct {
	mode   Mode
	length int
	// Maps shortPath to targetURL
	targetURLStore store.KVStore
	// Maps targetURL to shortPath
	shortPathStore store.KVStore
}

func NewURLShortner(length int, targetURLStore store.KVStore, shortPathStore store.KVStore) URLShortner {
	return URLShortner{
		mode:           Random,
		length:         length,
		targetURLStore: targetURLStore,
		shortPathStore: shortPathStore,
	}
}

func (u *URLShortner) GetTargetURL(shortPath string) (string, error) {
	targetURL, found, err := u.lookupTargetURL(shortPath)
	if err != nil {
		return "", ErrServerError
	}
	if !found {
		return "", store.ErrKeyNotFound
	}
	return targetURL, nil
}

func (u *URLShortner) CreateShortPath(shortPath string, targetURL string) (string, error) {
	if !isValidShortPath(shortPath) {
		return "", NewErrValidationFailed("shortPath is invalid")
	}
	if !isValidTargetURL(targetURL) {
		return "", NewErrValidationFailed("targetURL is invalid")
	}
	targetURL = removeTrailingSlash(targetURL)
	if len(shortPath) > 0 { // isShortPathProvidedInRequest ?
		oldTargetURL, found, err := u.lookupTargetURL(shortPath)
		if err != nil {
			return "", ErrServerError
		}
		if found {
			if oldTargetURL == targetURL {
				return shortPath, nil
			} else {
				return "", ErrConflict
			}
		}
	}
	existingShortPath, found, err := u.lookupShortPath(targetURL)
	if err != nil {
		return "", ErrServerError
	}
	if found { // isTargetURLAlreadyShortened
		return existingShortPath, nil
	}
	return u.doShorten(shortPath, targetURL)
}

func (u *URLShortner) doShorten(shortPath string, targetURL string) (string, error) {
	if len(shortPath) == 0 {
		shortPath = u.generateShortPath()
	}
	// FIXME: Get/Exists and Put calls from this file are not atomic.
	// This can lead to following inconsistent states.
	// i. existing shortpath gets replaced, both returns success
	// ii. Same targetURL gets shortened twice with different shortpaths
	err := u.targetURLStore.Put(shortPath, targetURL)
	if err != nil {
		return "", ErrServerError
	}
	err = u.shortPathStore.Put(targetURL, shortPath)
	if err != nil {
		u.targetURLStore.Delete(shortPath)
		return "", ErrServerError
	}
	return shortPath, nil
}

func (u *URLShortner) lookupTargetURL(shortPath string) (string, bool, error) {
	targetURL, err := u.targetURLStore.Get(shortPath)
	if err != nil {
		if err == store.ErrKeyNotFound {
			return "", false, nil
		}
		return "", false, ErrServerError
	}
	return targetURL, true, nil
}

func (u *URLShortner) lookupShortPath(targetURL string) (string, bool, error) {
	shortPath, err := u.shortPathStore.Get(targetURL)
	if err != nil {
		if err == store.ErrKeyNotFound {
			return "", false, nil
		}
		return "", false, ErrServerError
	}
	return shortPath, true, nil
}

func (u *URLShortner) generateShortPath() string {
	if u.mode == Phrase {
		// TODO: Implement
		panic("Phrase mode not implemented")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomBytes := make([]byte, u.length)
	for i := range randomBytes {
		randomBytes[i] = charset[r.Intn(len(charset))]
	}
	return string(randomBytes)
}

func isValidTargetURL(targetURL string) bool {
	return len(targetURL) == len(strings.TrimSpace(targetURL))
}

func isValidShortPath(shortPath string) bool {
	if len(shortPath) != len(strings.TrimSpace(shortPath)) {
		return false
	}
	for _, c := range shortPath {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
			return false
		}
	}
	return true
}

// removeTrailingSlash removes slash from the end of the targetURL
// this is because both URLs refer to the same resource
func removeTrailingSlash(targetURL string) string {
	if targetURL[len(targetURL)-1] == '/' {
		targetURL = targetURL[:len(targetURL)-1]
	}
	return targetURL
}

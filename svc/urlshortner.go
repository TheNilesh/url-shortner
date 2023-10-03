package svc

import (
	"context"
	"fmt"
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

func (u *URLShortner) GetTargetURL(ctx context.Context, shortPath string) (string, error) {
	targetURL, found, err := u.lookupTargetURL(shortPath)
	if err != nil {
		return "", NewErrServerError("could not lookup shortpath", err)
	}
	if !found {
		return "", NewErrNotFound("shortpath mapping not found")
	}
	return targetURL, nil
}

func (u *URLShortner) CreateShortPath(ctx context.Context, shortPath string, targetURL string) (string, error) {
	if err := validateShortPath(shortPath); err != nil {
		return "", err
	}
	if err := validateTargetURL(targetURL); err != nil {
		return "", err
	}
	targetURL = removeTrailingSlash(targetURL)
	if len(shortPath) > 0 { // isShortPathProvidedInRequest ?
		oldTargetURL, found, err := u.lookupTargetURL(shortPath)
		if err != nil {
			return "", fmt.Errorf("could not lookup shortpath: %w", err)
		}
		if found {
			if oldTargetURL == targetURL {
				return shortPath, nil
			} else {
				return "", NewErrConflict("shortpath already exists for different targetURL")
			}
		}
	}
	existingShortPath, found, err := u.lookupShortPath(targetURL)
	if err != nil {
		return "", err
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
		return "", NewErrServerError("could not save shortpath", err)
	}
	err = u.shortPathStore.Put(targetURL, shortPath)
	if err != nil {
		errDelete := u.targetURLStore.Delete(shortPath)
		if errDelete != nil {
			// TODO: Log error
			return "", NewErrServerError("could not delete shortpath from store", errDelete)
		}
		return "", NewErrServerError("could not save targetURL", err)
	}
	return shortPath, nil
}

func (u *URLShortner) lookupTargetURL(shortPath string) (string, bool, error) {
	targetURL, err := u.targetURLStore.Get(shortPath)
	if err != nil {
		if err == store.ErrKeyNotFound {
			return "", false, nil
		}
		return "", false, NewErrServerError("could not lookup shortpath for target URL", err)
	}
	return targetURL, true, nil
}

func (u *URLShortner) lookupShortPath(targetURL string) (string, bool, error) {
	shortPath, err := u.shortPathStore.Get(targetURL)
	if err != nil {
		if err == store.ErrKeyNotFound {
			return "", false, nil
		}
		return "", false, NewErrServerError("could not lookup shortpath for target URL", err)
	}
	return shortPath, true, nil
}

func (u *URLShortner) generateShortPath() string {
	if u.mode == Phrase {
		// TODO: Implement
		panic("phrase mode not implemented")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomBytes := make([]byte, u.length)
	for i := range randomBytes {
		randomBytes[i] = charset[r.Intn(len(charset))]
	}
	return string(randomBytes)
}

func validateTargetURL(targetURL string) error {
	if len(targetURL) != len(strings.TrimSpace(targetURL)) {
		return NewErrValidation("targetURL contains leading or trailing spaces")
	}
	return nil
}

func validateShortPath(shortPath string) error {
	if len(shortPath) != len(strings.TrimSpace(shortPath)) {
		return NewErrValidation("shortPath contains leading or trailing spaces")
	}
	for _, c := range shortPath {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
			return NewErrValidation("shortPath contains invalid characters")
		}
	}
	return nil
}

// removeTrailingSlash removes slash from the end of the targetURL
// this is because both URLs refer to the same resource
func removeTrailingSlash(targetURL string) string {
	if targetURL[len(targetURL)-1] == '/' {
		targetURL = targetURL[:len(targetURL)-1]
	}
	return targetURL
}

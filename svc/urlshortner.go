package svc

import (
	"context"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/thenilesh/url-shortner/metrics"
	"github.com/thenilesh/url-shortner/store"
)

type Mode string

const (
	Random Mode = "Random"
	Phrase Mode = "Phrase"

	charset = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type URLShortner interface {
	GetTargetURL(ctx context.Context, shortPath string) (string, error)
	CreateShortPath(ctx context.Context, shortPath string, targetURL string) (string, error)
}

type urlShortner struct {
	mode   Mode
	length int
	// Maps shortPath to targetURL
	targetURLStore store.KVStore
	// Maps targetURL to shortPath
	shortPathStore store.KVStore
	metrics        metrics.Metrics
}

func NewURLShortner(length int, targetURLStore store.KVStore, shortPathStore store.KVStore, metrics metrics.Metrics) URLShortner {
	return &urlShortner{
		mode:           Random,
		length:         length,
		targetURLStore: targetURLStore,
		shortPathStore: shortPathStore,
		metrics:        metrics,
	}
}

func (u *urlShortner) GetTargetURL(ctx context.Context, shortPath string) (string, error) {
	targetURL, found, err := u.lookupTargetURL(ctx, shortPath)
	if err != nil {
		return "", NewErrServerError("could not lookup shortpath", err)
	}
	if !found {
		return "", NewErrNotFound("shortpath mapping not found")
	}
	return targetURL, nil
}

func (u *urlShortner) CreateShortPath(ctx context.Context, shortPath string, targetURL string) (string, error) {
	if err := validateShortPath(shortPath); err != nil {
		return "", err
	}
	if err := validateTargetURL(targetURL); err != nil {
		return "", err
	}
	targetURL = removeTrailingSlash(targetURL)
	if len(shortPath) > 0 { // isShortPathProvidedInRequest ?
		oldTargetURL, found, err := u.lookupTargetURL(ctx, shortPath)
		if err != nil {
			return "", NewErrServerError("could not lookup shortpath", err)
		}
		if found {
			if oldTargetURL == targetURL {
				return shortPath, nil
			} else {
				return "", NewErrConflict("shortpath already exists for different targetURL")
			}
		}
	}
	existingShortPath, found, err := u.lookupShortPath(ctx, targetURL)
	if err != nil {
		return "", err
	}
	if found { // isTargetURLAlreadyShortened
		return existingShortPath, nil
	}
	return u.doShorten(ctx, shortPath, targetURL)
}

func (u *urlShortner) doShorten(ctx context.Context, shortPath string, targetURL string) (string, error) {
	var shouldGenerateShortPath bool
	if len(shortPath) == 0 {
		shortPath = u.generateShortPath()
		shouldGenerateShortPath = true
	}
	// FIXME: Get/Exists and Put calls from this file are not atomic.
	// This can lead to following inconsistent states.
	// i. existing shortpath gets replaced, both returns success
	// ii. Same targetURL gets shortened twice with different shortpaths
	var err error
	var found bool
	for i := 0; i < 3; i++ {
		if found, err = u.targetURLStore.Exists(ctx, shortPath); err != nil {
			return "", NewErrServerError("could not check if shortpath exists", err)
		}
		if found && shouldGenerateShortPath {
			shortPath = u.generateShortPath()
			continue
		}
		if found {
			return "", NewErrConflict("shortpath already exists")
		}
		err = u.targetURLStore.Put(ctx, shortPath, targetURL)
		if err != nil {
			return "", NewErrServerError("could not save shortpath", err)
		}
		err = u.shortPathStore.Put(ctx, targetURL, shortPath)
		if err != nil {
			errDelete := u.targetURLStore.Delete(ctx, shortPath)
			if errDelete != nil {
				// TODO: Log error
				return "", NewErrServerError("could not delete shortpath from store", errDelete)
			}
			return "", NewErrServerError("could not save targetURL", err)
		}
		u.metrics.GetCollector("domain_shortens").Inc(extractDomainFromURL(targetURL))
		return shortPath, nil
	}
	return "", NewErrServerError("could not find available short path", nil)
}

func (u *urlShortner) lookupTargetURL(ctx context.Context, shortPath string) (string, bool, error) {
	targetURL, err := u.targetURLStore.Get(ctx, shortPath)
	if err != nil {
		if err == store.ErrKeyNotFound {
			return "", false, nil
		}
		return "", false, NewErrServerError("could not lookup shortpath for target URL", err)
	}
	return targetURL, true, nil
}

func (u *urlShortner) lookupShortPath(ctx context.Context, targetURL string) (string, bool, error) {
	shortPath, err := u.shortPathStore.Get(ctx, targetURL)
	if err != nil {
		if err == store.ErrKeyNotFound {
			return "", false, nil
		}
		return "", false, NewErrServerError("could not lookup shortpath for target URL", err)
	}
	return shortPath, true, nil
}

func (u *urlShortner) generateShortPath() string {
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
	_, err := url.Parse(targetURL)
	if err != nil {
		return NewErrValidation("targetURL is not valid")
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
	if shortPath == "metrics" {
		return NewErrValidation("shortPath is reserved")
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

func extractDomainFromURL(rawURL string) string {
	u, _ := url.Parse(rawURL)
	return u.Hostname()
}

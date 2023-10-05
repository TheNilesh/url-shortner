package svc

import (
	"context"
	"net/url"
	"strings"

	"github.com/thenilesh/url-shortner/metrics"
	"github.com/thenilesh/url-shortner/store"
)

type URLShortner interface {
	GetTargetURL(ctx context.Context, shortPath string) (string, error)
	CreateShortPath(ctx context.Context, shortPath string, targetURL string) (string, error)
}

type urlShortner struct {
	randomStrGen RandomStrGen
	// Maps shortPath to targetURL
	targetURLStore store.KVStore
	// Maps targetURL to shortPath
	shortPathStore store.KVStore
	metrics        metrics.Metrics
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
	if len(shortPath) == 0 {
		shortPath, err = u.findAvailableShortPath(ctx, u.targetURLStore)
		if err != nil {
			return "", err
		}
	}
	return u.doShorten(ctx, shortPath, targetURL)
}

// findAvailableShortPath finds a randomly generated shortPath that is not already taken
func (u *urlShortner) findAvailableShortPath(ctx context.Context, store store.KVStore) (string, error) {
	var err error
	var alreadyExists bool
	for i := 0; i < 3; i++ {
		shortPath := u.randomStrGen.Generate()
		if alreadyExists, err = u.targetURLStore.Exists(ctx, shortPath); err != nil {
			return "", NewErrServerError("could not check if shortpath exists", err)
		}
		if !alreadyExists {
			return shortPath, nil
		}
	}
	return "", NewErrServerError("failed to generate available short_path", nil)
}

func (u *urlShortner) doShorten(ctx context.Context, shortPath string, targetURL string) (string, error) {

	// FIXME: Get/Exists and Put calls from this file are not atomic.
	// This can lead to following inconsistent states.
	// i. existing shortpath gets replaced, both returns success
	// ii. Same targetURL gets shortened twice with different shortpaths

	err := u.targetURLStore.Put(ctx, shortPath, targetURL)
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
		return NewErrValidation("short_path contains leading or trailing spaces")
	}
	for _, c := range shortPath {
		// Allow alphanumeric, - and _
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return NewErrValidation("short_path contains disallowed characters")
		}
	}
	if shortPath == "metrics" {
		return NewErrValidation("this short_path is reserved")
	}
	if len(shortPath) > 50 {
		return NewErrValidation("short_path is too long")
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

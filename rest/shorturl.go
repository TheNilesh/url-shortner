package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/thenilesh/url-shortner/store"
	"github.com/thenilesh/url-shortner/svc"
)

type ShortURL struct {
	store       store.KVStore
	revStore    store.KVStore
	urlShortner svc.URLShortner
}

func NewShortURL(store store.KVStore, revStore store.KVStore, urlShortner svc.URLShortner) *ShortURL {
	return &ShortURL{
		store:       store,
		revStore:    revStore,
		urlShortner: urlShortner,
	}
}

type URL struct {
	ShortPath string `json:"short_path"` // User provided short path
	TargetURL string `json:"target_url"`
}

func (s *ShortURL) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: This function is too big, break it down into smaller functions
	// To break-down, lets list down responsibilities of this function
	// 1. Decode JSON request body
	// 2. Check if shortPath is already taken
	// 3. Check if targetURL is already shortened
	// 4. If targetURL is already shortened, check if shortPath is same as user provided
	// 5. If shortPath is not provided, generate a new shortPath
	// 6. If shortPath is provided, check if it is already taken
	// 7. If shortPath is already taken, return 409
	// 8. If shortPath is not taken, store the shortPath and targetURL in store
	// 9. If storing shortPath and targetURL is successful, return 201
	// 10. If storing shortPath and targetURL is not successful, return 500

	// TODO: Prevent OOM/buffer overflow by not parsing large request body
	var req URL
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	// TODO: Validate request body, shortPath should be alphanumeric, no spaces allowed

	shortURL, err := s.revStore.Get(req.TargetURL)
	if err == nil {
		if req.ShortPath != "" && req.ShortPath != shortURL {
			w.WriteHeader(http.StatusConflict)
			return
		}
		shortURL = fmt.Sprintf("/%s", shortURL)
		w.Header().Set("Location", fmt.Sprintf("/%s", shortURL))
		w.WriteHeader(http.StatusCreated)
		return
	}
	if err != nil && err != store.ErrKeyNotFound {
		// TODO: Log
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Cache backend may not be fast, consider using context
	// and map it with http timeout. ie, if http request times out,
	// cancel the context

	if req.ShortPath != "" {
		targetURL, err := s.store.Get(req.ShortPath)
		if err != nil && err != store.ErrKeyNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err == nil {
			if targetURL == req.TargetURL {
				shortURL = fmt.Sprintf("/%s", req.ShortPath)
				w.Header().Set("Location", fmt.Sprintf("/%s", shortURL))
				w.WriteHeader(http.StatusCreated)
				return
			} else {
				w.WriteHeader(http.StatusConflict)
				return
			}
		}
	}

	if req.ShortPath == "" {
		// FIXME: Following loop will execute indefinitely if shortPaths are exhausted
		for exists := true; exists; {
			shortURL = s.urlShortner.Shorten(req.TargetURL)
			exists, err = s.store.Exists(shortURL)
			if err != nil {
				// TODO: Log
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	} else {
		shortURL = req.ShortPath
	}

	err = s.store.Put(shortURL, req.TargetURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = s.revStore.Put(req.TargetURL, shortURL)
	if err != nil {
		// TODO: Log
		err := s.store.Delete(shortURL)
		if err != nil {
			// TODO: Log
			fmt.Println("Failed to delete shortURL", err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/%s", shortURL))
	w.WriteHeader(http.StatusCreated)
}

func (s *ShortURL) Get(w http.ResponseWriter, r *http.Request) {
	// TODO: if do not redirect query param is included then return ShortURL resource
	// otherwise redirect user to the targetURL
	vars := mux.Vars(r)
	shortPath := vars["id"]
	targetURL, err := s.store.Get(shortPath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.Redirect(w, r, targetURL, http.StatusMovedPermanently)
}

func (s *ShortURL) Put(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}

func (s *ShortURL) Delete(w http.ResponseWriter, r *http.Request) {
	//TODO: Implement
}
